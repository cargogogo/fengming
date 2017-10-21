package agent

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/cargogogo/fengming/model"
	_ "github.com/cargogogo/fengming/utils/loghook"
	"github.com/cargogogo/fengming/utils/reqlog"
)

// Load loads the router
func Load(cfg *model.AgentConfig) http.Handler {

	logrus.Debugf("\n\nLoad with config:\n %+v\n\n", cfg)

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(reqlog.ReqLoggerMiddleware(logrus.New(), time.RFC3339, true))

	svc := AgentService{cfg}

	v1group := e.Group("/v1")
	{
		v1group.POST("/task", svc.PostTask)
	}

	return e
}
