package common

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	CmdCopy          = "cp"
	CmdTorrent       = "torrent"
	CmdTorrentCreate = "torrent-create"
)

func CreateTorrentFile(trackerAddr, dataPath, torrentFile string) error {
	ctx := context.TODO()

	cmds := []string{
		CmdCopy,
		"-r",
		dataPath,
		".",
	}
	log.Infof("cmds: %v", cmds)
	if err := ExecCmdNoOutput(ctx, cmds); err != nil {
		return err
	}

	cmds = []string{
		CmdTorrentCreate,
		fmt.Sprintf("-a=udp://%s", trackerAddr),
		fmt.Sprintf("./%s", filepath.Base(dataPath)),
	}
	log.Infof("cmds: %v", cmds)

	output, err := ExecCmd(ctx, cmds)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(torrentFile, output, 0655)
}

func PullFromTorrent(torrentFile string, needSeed bool, listenAddr string) error {
	ctx := context.TODO()

	cmds := []string{CmdTorrent}
	if needSeed {
		cmds = append(cmds, "-seed=true", fmt.Sprintf("-addr=%s", listenAddr))
	}
	cmds = append(cmds, torrentFile)
	log.Infof("cmds: %v", cmds)

	return ExecCmdNoOutput(ctx, cmds)
}
