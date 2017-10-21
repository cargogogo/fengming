package model

import (
	"time"
)

// AgentConfig ...
type AgentConfig struct {
	MasterAddr     string
	NodeName       string
	DownloadDir    string
	ListenAddr     string
	ReportInterval time.Duration
}

// Task ...
type Task struct {
	ID          string
	LayerName   string
	Status      string
	Torrent     []byte
	TorrentPath string
}

// AgentStatus ...
type AgentStatus struct {
	Name           string
	Addr           string // ip:port
	APIPushTorrent string // POST http://ip:port/torrent
	Tasks          []Task
}
