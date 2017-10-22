package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ianschenck/envflag"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"

	"github.com/cargogogo/fengming/model"
	"github.com/cargogogo/fengming/pkg/agent"
)

// ServeCommand exports the server command.
var ServeCommand = cli.Command{
	Name:   "serve",
	Usage:  "starts the agent daemon",
	Action: server,
	Flags: []cli.Flag{
		cli.BoolFlag{
			EnvVar: "DEBUG",
			Name:   "debug",
			Usage:  "start the agent in debug mode",
		},
		cli.StringFlag{
			EnvVar: "SERVER_ADDR",
			Name:   "server-addr",
			Usage:  "server address",
			Value:  ":7100",
		},
		cli.StringFlag{
			EnvVar: "MASTER_ADDR",
			Name:   "master-addr",
			Usage:  "master addr",
			Value:  "http://controller:7000",
		},
		cli.StringFlag{
			EnvVar: "NODE_NAME",
			Name:   "node-name",
			Usage:  "node name",
		},
		cli.DurationFlag{
			EnvVar: "REPORT_INTERVAL",
			Name:   "report-interval",
			Usage:  "node name",
			Value:  time.Second * 5,
		},
		cli.StringFlag{
			EnvVar: "DOWN_DIR",
			Name:   "download-dir",
			Usage:  "download dir",
			Value:  "/tmp",
		},
	},
}

func server(c *cli.Context) error {
	// debug level if requested by user
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
	cfg := &model.AgentConfig{}
	cfg.MasterAddr = c.String("master-addr")
	cfg.NodeName = c.String("node-name")
	cfg.ReportInterval = c.Duration("report-interval")
	cfg.DownloadDir = c.String("download-dir")
	cfg.ListenAddr = c.String("server-addr")

	// setup the server and start the listener
	handler := agent.Load(cfg)

	// start the server
	return http.ListenAndServe(
		c.String("server-addr"),
		handler,
	)
}

func main() {
	envflag.Parse()

	app := cli.NewApp()
	app.Name = "fengming"
	app.Version = "1.0"
	app.Usage = "command line utility"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		ServeCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
