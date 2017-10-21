package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	TrackerAddr = "45.76.163.62:6881"

	CurrentDir, _ = os.Getwd()
	DataPath      = filepath.Join(CurrentDir, "bt.go")
	TorrentFile   = filepath.Join(CurrentDir, "bt.torrent")
)

func TestCreateTorrentFile(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(CreateTorrentFile(TrackerAddr, DataPath, TorrentFile))
	defer func() {
		assert.NoError(os.Remove(TorrentFile))
	}()
}
