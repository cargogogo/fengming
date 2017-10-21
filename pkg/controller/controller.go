package controller

import (
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

type ControllerI interface {
	// CreateTorrent creates a torrent file from raw data path.
	CreateTorrent(rawDataPath, targetTorrent string) error

	// DistrbuteTorrent distribute a torrent file to the hosts.
	DistributeTorrent(torrent string, hosts []string) error
}

type ControllerConfig struct {
	AnnounceList [][]string
	PieceLength  int64

	RootPath string
}

type controller struct {
	config *ControllerConfig
}

func NewController(cfg *ControllerConfig) (ControllerI, error) {
	return &controller{
		config: cfg,
	}, nil
}

func (c *controller) CreateTorrent(rawDataPath, targetTorrent string) error {
	mi := metainfo.MetaInfo{
		AnnounceList: c.config.AnnounceList,
	}
	mi.SetDefaults()

	info := metainfo.Info{
		PieceLength: c.config.PieceLength,
	}
	if err := info.BuildFromFilePath(c.config.RootPath); err != nil {
		return err
	}

	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		return err
	}

	return nil
}

func (c *controller) DistributeTorrent(torrent string, hosts []string) error {
	return nil
}
