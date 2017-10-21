package controller

import (
	"github.com/cargogogo/fengming/pkg/common"
)

type ControllerI interface {
	// CreateTorrent creates a torrent file from raw data path.
	CreateTorrent(rawDataPath, targetTorrent string) error

	// SeedTorrent seeds the data.
	SeedTorrent(torrentFile, listenAddr string) error

	// DistrbuteTorrent distribute a torrent file to the hosts.
	DistributeTorrent(torrent string, hosts []string) error
}

type ControllerConfig struct {
	TrackerAddr string
}

type controller struct {
	config *ControllerConfig
}

func NewController(cfg *ControllerConfig) (ControllerI, error) {
	return &controller{
		config: cfg,
	}, nil
}

func (c *controller) CreateTorrent(dataPath, torrentFile string) error {
	return common.CreateTorrentFile(c.config.TrackerAddr, dataPath, torrentFile)
}

func (c *controller) SeedTorrent(torrentFile, listenAddr string) error {
	return common.PullFromTorrent(torrentFile, true, listenAddr)
}

func (c *controller) DistributeTorrent(torrent string, hosts []string) error {
	return nil
}
