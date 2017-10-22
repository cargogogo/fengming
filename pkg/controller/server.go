package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/cargogogo/fengming/model"
	httpClient "github.com/cargogogo/fengming/utils/http"
)

const ActionPush = "push"

const (
	MinPort = 20000
	MaxPort = 30000
)

type ServerConfig struct {
	Addr string

	// abs path, similar to: /home/ubuntu/registry/docker/registry/v2/blobs/sha256/
	RegistryBlobPath string

	TrackerAddr string

	Logger *log.Entry
}

type Server struct {
	config *ServerConfig

	router *gin.Engine

	controller ControllerI

	logger *log.Entry

	agents map[string]model.AgentStatus

	filter model.Filter

	client httpClient.Client
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = log.WithFields(log.Fields{
			"app": "server",
		})
	}

	controller, err := NewController(&ControllerConfig{
		TrackerAddr: cfg.TrackerAddr,
	})
	if err != nil {
		return nil, err
	}

	return &Server{
		config:     cfg,
		router:     gin.Default(),
		controller: controller,
		agents:     map[string]model.AgentStatus{},
		logger:     logger,
		client:     httpClient.DefaultClient,
	}, nil
}

func (s *Server) Run() error {
	s.router.POST("/v1/hook", s.RegistryHook)
	s.router.POST("/v1/agents", s.AgentHeartbeat)
	s.router.GET("/v1/agents", s.AgentsInfo)

	s.router.GET("/v1/filter", s.GetFilter)
	s.router.POST("/v1/filter", s.PostFilter)

	return s.router.Run(s.config.Addr)
}

func (s *Server) RegistryHook(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	// Make and distribute torrent files for `PUSH` action.
	var data Events
	if err := decoder.Decode(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("failed to unmarshal data: %s", err),
		})
		return
	}
	if data.Events[0].Action != ActionPush {
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("got action: %s", data.Events[0].Action),
		})
		return
	}

	// Get the absolute path of the layer data.
	digest := data.Events[0].Target.Digest[7:]
	dataDir := filepath.Join(s.config.RegistryBlobPath, digest[:2], digest)
	s.logger.Infof("dataDir: %s", dataDir)

	// Make a torrent file for one layer of docker image.
	// Note: the torrent file should be with the data dir.
	torrentFile := digest + ".torrent"
	if err := s.controller.CreateTorrent(dataDir, torrentFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("failed to make torrent file: %s", err),
		})
		return
	}

	// Seed the torrent.
	// TODO: It's better to assign the port by a service other than using random
	// port.
	port := MinPort + rand.Intn(MaxPort-MinPort)
	seedListenAddr := fmt.Sprintf("%s:%d", strings.Split(s.config.Addr, ":")[0], port)
	s.logger.Infof("listen addr: %s", seedListenAddr)
	s.SeedTorrent(torrentFile, seedListenAddr)

	// TODO: Distribute the torrent file.
	dataBytes, err := ioutil.ReadFile(torrentFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("failed to read torrent data: %s", err),
		})
		return
	}

	// TODO: Do this in the background.
	ctx := context.TODO()
	for _, agent := range s.agents {
		requestBody := model.Task{
			LayerName: digest,
			Torrent:   dataBytes,
		}
		var ret interface{}
		if err := s.client.CallWithJson(ctx, ret, "POST", agent.APIPushTorrent, requestBody); err != nil {
			s.logger.Warnf("Failed to post torrent %s to %s: %s",
				torrentFile, agent.Name, err)
		} else {
			s.logger.Infof("Successfully post torrent %s to %s", torrentFile, agent.Name)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "register hook",
	})
}

func (s *Server) AgentHeartbeat(c *gin.Context) {
	var agentStatus model.AgentStatus
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&agentStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("Failed to unmarshal agent data: %s", err),
		})
		return
	}
	s.logger.Infof("Receive agent heartbeat: %v", agentStatus)

	s.agents[agentStatus.Name] = agentStatus

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok from AgentHeartbeat",
	})
}

func (s *Server) AgentsInfo(c *gin.Context) {
	status := []model.AgentStatus{}
	for _, agent := range s.agents {
		status = append(status, agent)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "ok from AgentsInfo",
		"status": status,
	})
}

// TODO: Add timeout.
func (s *Server) SeedTorrent(torrentFile, listenAddr string) {
	go s.controller.SeedTorrent(torrentFile, listenAddr)
}
