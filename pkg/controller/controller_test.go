package controller

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cargogogo/fengming/pkg/common"
)

func TestCreateTorrent(t *testing.T) {
	assert := assert.New(t)

	cfg := &ControllerConfig{
		AnnounceList: [][]string{
			[]string{"udp://test.com:80"},
		},
		PieceLength: 1024 * 1024 * 4,
	}
	controller, err := NewController(cfg)
	assert.NoError(err)

	dataDir := "./tmp"
	os.MkdirAll(dataDir, 0766)
	defer func() {
		os.RemoveAll(dataDir)
	}()

	ctx := context.TODO()
	cmds := []string{
		"cp",
		filepath.Join(os.Getenv("HOME"), "go1.9.linux-amd64.tar.gz"),
		dataDir,
	}
	assert.NoError(common.ExecCmdNoOutput(ctx, cmds))

	torrentFile := "go1.9.torrent"
	assert.NoError(controller.CreateTorrent(dataDir, torrentFile))
	defer func() {
		os.Remove(torrentFile)
	}()
}
