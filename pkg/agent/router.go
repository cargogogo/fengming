package agent

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cargogogo/fengming/model"
	"github.com/cargogogo/fengming/utils/header"
	_ "github.com/cargogogo/fengming/utils/loghook"
	"github.com/cargogogo/fengming/utils/reqlog"
)

// Load loads the router
func Load(cfg *model.AgentConfig) http.Handler {

	logrus.Debugf("\n\nLoad with config:\n %+v\n\n", cfg)

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.Secure)
	e.Use(header.Options)

	e.Use(reqlog.ReqLoggerMiddleware(logrus.New(), time.RFC3339, true))

	svc := New(cfg)
	go svc.Run()

	v1group := e.Group("/v1")
	{
		v1group.POST("/task", svc.PostTask)
		v1group.GET("/test", test)
		v1group.GET("/testfilter", testfilter)
		v1group.POST("/testfilter", testfilter)
	}

	return e
}

func test(c *gin.Context) {
	agents := []model.AgentStatus{}
	for i := 1; i < 3; i++ {
		status := model.AgentStatus{
			Name: "123",
			Addr: "1.1.1.1:2222",
		}
		for i := 1; i < 3; i++ {
			status.Tasks = append(status.Tasks, model.Task{
				ID:        "12123",
				LayerName: "12w123213",
				Status:    "complete",
			})
		}
		agents = append(agents, status)
	}
	c.JSON(200, agents)
}

func testfilter(c *gin.Context) {
	filter := &model.Filter{
		AgentName: "123",
		Repo:      "456",
	}
	c.JSON(200, filter)
}
