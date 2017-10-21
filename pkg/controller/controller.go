package controller

import (
	"os"
	"time"

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
	AnnounceList [][]string // tracker addresses
	PieceLength  int64      // data chunk size
}

type controller struct {
	config *ControllerConfig
}

func NewController(cfg *ControllerConfig) (ControllerI, error) {
	return &controller{
		config: cfg,
	}, nil
}

func (c *controller) CreateTorrent(dataDir, targetTorrent string) error {
	mi := metainfo.MetaInfo{
		AnnounceList: c.config.AnnounceList,
		Comment:      "Distribute images based on P2P",
		CreatedBy:    "FengMing",
		CreationDate: time.Now().Unix(),
	}

	info := metainfo.Info{
		PieceLength: c.config.PieceLength,
	}
	if err := info.BuildFromFilePath(dataDir); err != nil {
		return err
	}

	var err error
	if mi.InfoBytes, err = bencode.Marshal(info); err != nil {
		return err
	}

	fp, err := os.Create(targetTorrent)
	if err != nil {
		return err
	}
	defer fp.Close()

	return mi.Write(fp)
}

func (c *controller) DistributeTorrent(torrent string, hosts []string) error {
	return nil
}
