package agent

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cargogogo/fengming/model"

	"github.com/dustin/go-humanize"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func (d *layerloader) updatestatus(t *torrent.Torrent) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Minute * 10) // sleep one second
		timeout <- true
	}()
	for {
		select {
		case <-t.GotInfo():
		case <-timeout:
			return
		}
		if t.Seeding() {
			d.task.Status = "seeding"
			d.upfunc(&d.task)
		} else if t.BytesCompleted() == t.Info().TotalLength() {
			d.task.Status = "completed"
			d.upfunc(&d.task)
		} else {
			d.task.Status = fmt.Sprintf("downloading (%s/%s)", humanize.Bytes(uint64(t.BytesCompleted())), humanize.Bytes(uint64(t.Info().TotalLength())))
			d.upfunc(&d.task)
		}
	}
}

func (d *layerloader) addTorrents(client *torrent.Client) (err error) {
	var t *torrent.Torrent
	if strings.HasPrefix(d.task.TorrentPath, "magnet:") {
		t, err = client.AddMagnet(d.task.TorrentPath)
		if err != nil {
			return
		}
	} else {
		var meta *metainfo.MetaInfo
		reader := bytes.NewReader(d.task.Torrent)
		meta, err = metainfo.Load(reader)
		if err != nil {
			return
		}
		t, err = client.AddTorrent(meta)
		if err != nil {
			return
		}
	}

	go func() {
		<-t.GotInfo()
		t.DownloadAll()
	}()
	go d.updatestatus(t)
	return nil
}

type updatefunc func(*model.Task)
type layerloader struct {
	task   model.Task
	cfg    torrent.Config
	upfunc updatefunc
}

func (d *layerloader) load() error {
	client, err := torrent.NewClient(&d.cfg)
	if err != nil {
		return err
	}
	defer func() {
		go func() {
			time.Sleep(time.Minute * 9)
			client.Close()
		}()

		d.task.Status = "finish"
		d.upfunc(&d.task)
	}()

	err = d.addTorrents(client)
	if client.WaitAll() {
		logrus.Debug("compelte task ", d.task)
		return nil
	}
	return errors.New("error not sure")
}
