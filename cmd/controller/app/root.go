package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cargogogo/fengming/pkg/controller"
)

var (
	debugLevel       uint32
	addr             string
	registryBlobPath string
	trackerAddr      string
)

const (
	defaultDebugLevel uint32 = 4
	defaultAddr              = "0.0.0.0:9090"
)

var rootCmd = &cobra.Command{
	Use:          "controller",
	Short:        "controller provides a way to make and distribute torrent files",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.WithFields(log.Fields{
			"app": "controller",
		})

		config := &controller.ServerConfig{
			Addr:             addr,
			RegistryBlobPath: registryBlobPath,
			TrackerAddr:      trackerAddr,
			Logger:           logger,
		}
		server, err := controller.NewServer(config)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(server.Run())
	},
}

func init() {
	rootCmd.PersistentFlags().Uint32VarP(&debugLevel, "debuglevel", "l", defaultDebugLevel,
		"log debug level: 0[panic] 1[fatal] 2[error] 3[warn] 4[info] 5[debug]")
	rootCmd.PersistentFlags().StringVar(&addr, "addr", defaultAddr, "listen address")
	rootCmd.PersistentFlags().StringVar(&registryBlobPath, "registryBlobPath",
		"", "path to the blob dir")
	rootCmd.PersistentFlags().StringVar(&trackerAddr, "trackerAddr",
		"", "the address of the tracker")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// set logging
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.Level(debugLevel))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
