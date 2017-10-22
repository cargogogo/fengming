package controller

import (
	"regexp"

	"github.com/cargogogo/fengming/model"
	"github.com/gin-gonic/gin"
)

// PostFilter ...
func (s *Server) PostFilter(c *gin.Context) {
	err := c.BindJSON(&s.filter)
	if err != nil {
		c.AbortWithError(400, err)
	}
	c.JSON(200, s.filter)
}

// GetFilter ...
func (s *Server) GetFilter(c *gin.Context) {
	c.JSON(200, s.filter)
}

func (s *Server) filterEvent(events []Event) (ret []Event) {
	if s.filter.Repo == "" {
		return events
	}

	for _, e := range events {
		match, err := regexp.MatchString(s.filter.Repo, e.Target.Repository)
		if err != nil || match {
			ret = append(ret, e)
		}
	}
	return
}

func (s *Server) filterAgent(agents []model.AgentStatus) (ret []model.AgentStatus) {
	if s.filter.AgentName == "" {
		return agents
	}

	for _, e := range agents {
		match, err := regexp.MatchString(s.filter.AgentName, e.Name)
		if err != nil || match {
			ret = append(ret, e)
		}
	}
	return
}
