package agent

import (
	"github.com/cargogogo/fengming/model"

	"github.com/gin-gonic/gin"
)

// AgentService ...
type AgentService struct {
	*model.AgentConfig
}

// PostTask ...
func (s *AgentService) PostTask(c *gin.Context) {
	c.AbortWithStatus(200)
}
