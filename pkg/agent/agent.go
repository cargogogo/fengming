package agent

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cargogogo/fengming/model"
	httpclient "github.com/cargogogo/fengming/utils/http"

	"github.com/anacrolix/dht"
	"github.com/anacrolix/torrent"
	"github.com/cargogogo/fengming/pkg/common"
	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var constPortStart = 21000
var portStart = constPortStart
var constPortEnd = 22000
var portLock sync.RWMutex

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Service ...
type Service struct {
	*model.AgentConfig
	tasks     chan *model.Task
	taskqueue map[string]*model.Task
	tasklock  sync.RWMutex
	// docker    *client.Client
}

// New ....
func New(cfg *model.AgentConfig) *Service {
	// cli, err := client.NewEnvClient()
	// if err != nil {
	// 	logrus.Fatalln(err)
	// }
	return &Service{
		AgentConfig: cfg,
		tasks:       make(chan *model.Task, 10),
		taskqueue:   make(map[string]*model.Task),
		// docker:      cli,
	}
}

// curl -X POST -d '{"torrentPath":"magnet:?xt=urn:btih:9a226c0dac90a56a148fc58eff5a38aa7adaddf7", "layerName", "dat1.zip"}' http://localhost:7100/v1/task

// PostTask ...
func (s *Service) PostTask(c *gin.Context) {
	c.AbortWithStatus(200)
	t := &model.Task{}
	err := c.BindJSON(t)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	t.Status = "downloading"
	// t.ID = randStringRunes(10)
	t.ID = t.LayerName

	s.tasks <- t
	c.JSON(200, gin.H{})
}

// Run report status to master
func (s *Service) Run() {
	tick := time.Tick(s.ReportInterval)
	logrus.Debugln("run with interval", s.ReportInterval)
	for {
		select {
		case <-tick:
			s.reportStatus()
		case task := <-s.tasks:
			s.addTask(task)
		}
	}
}

func (s *Service) reportStatus() {
	s.tasklock.RLock()
	taskcopy := []model.Task{}
	for _, t := range s.taskqueue {
		a := *t
		a.Torrent = []byte{}
		a.ID = t.ID
		taskcopy = append(taskcopy, a)
	}
	s.tasklock.RUnlock()
	status := &model.AgentStatus{
		Name:           s.NodeName,
		Addr:           s.ListenAddr,
		APIPushTorrent: "http://" + s.ListenAddr + "/v1/task",
		Tasks:          taskcopy,
	}
	logrus.Debugf("reportStatus: %d %+v\n", len(status.Tasks), status)

	err := httpclient.DefaultClient.CallWithJson(nil, nil, "POST", s.MasterAddr+"/v1/agents", status)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *Service) addTask(t *model.Task) {
	s.tasklock.Lock()
	defer s.tasklock.Unlock()
	s.taskqueue[t.ID] = t
	go s.taskrun(t)
}

func (s *Service) rmtask(t *model.Task) {
	s.tasklock.Lock()
	defer s.tasklock.Unlock()
	delete(s.taskqueue, t.ID)
}

func (s *Service) taskrun(t *model.Task) {
	var err error
	exsit := s.checkLayerExsit(t)
	if !exsit {
		err = s.downloadTask(t)
		logrus.Debug("downloadTask ", err)
		if err == nil {
			err = s.importLayer(t)
		}
	}
	if err != nil {
		logrus.Error(err)
		t.Status = "fail"
	}
}

func (s *Service) checkLayerExsit(t *model.Task) bool {
	return false
}

func (s *Service) updateTask(t *model.Task) {
	s.tasklock.Lock()
	defer s.tasklock.Unlock()
	if a, ok := s.taskqueue[t.ID]; ok {
		a.Status = t.Status
		if a.Status == "finish" {
			delete(s.taskqueue, t.ID)
		}
	}
}

func (s *Service) downloadTask(t *model.Task) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	portLock.Lock()
	portStart++
	if portStart == constPortEnd {
		portStart = constPortStart
	}
	portLock.Unlock()
	loader := layerloader{
		task:   *t,
		upfunc: s.updateTask,
		cfg: torrent.Config{
			DHTConfig: dht.ServerConfig{
				StartingNodes: dht.GlobalBootstrapAddrs,
			},
			Seed:       true,
			DataDir:    s.DownloadDir,
			Debug:      true,
			ListenAddr: ":" + strconv.Itoa(portStart),
		},
	}
	err := loader.load()
	return err
}

func (s *Service) importLayer(t *model.Task) error {
	ctx := context.Background()
	name := s.DownloadDir + "/" + t.LayerName + "/data"
	logrus.Debug("importing layer ", name)
	out, err := common.ExecCmd(ctx, []string{"docker", "import", name})
	logrus.Debug("importing layer ", name, string(out), err)
	return err
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
