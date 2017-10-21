package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const ActionPush = "push"

type AgentInfo struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

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

	agents []AgentInfo
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
		logger:     logger,
	}, nil
}

func (s *Server) Run() error {
	s.router.POST("/v1/hook", s.RegistryHook)
	s.router.POST("/v1/agents", s.AgentHeartbeat)
	s.router.GET("/v1/agents", s.AgentsInfo)

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

	// TODO: Make a torrent file for one layer of docker image.
	// Note: the torrent file should be with the data dir.
	torrentFile := digest + ".torrent"
	if err := s.controller.CreateTorrent(dataDir, torrentFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("failed to make torrent file: %s", err),
		})
		return
	}

	// TODO: Distribute the torrent file.

	c.JSON(http.StatusOK, gin.H{
		"msg": "register hook",
	})
}

func (s *Server) AgentHeartbeat(c *gin.Context) {
	var agentInfo AgentInfo
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&agentInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("Failed to unmarshal agent data: %s", err),
		})
		return
	}
	s.logger.Infof("Receive agent heartbeat: %v", agentInfo)

	s.agents = append(s.agents, agentInfo)

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok from AgentHeartbeat",
	})
}

func (s *Server) AgentsInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok from AgentsInfo",
	})
}
