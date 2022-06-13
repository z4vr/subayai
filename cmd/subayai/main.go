package main

import (
	"flag"
	"github.com/sarulabs/di/v2"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/discord"
	"github.com/z4vr/subayai/internal/services/slashcommands"
	"github.com/zekrotja/ken"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
)

var (
	flagConfigPath = flag.String("config", "config.toml", "Path to config file")
	flagCPUProfile = flag.String("cpuprofile", "", "Path to write CPU profile")
)

func main() {

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	if *flagCPUProfile != "" {
		f, err := os.Create(*flagCPUProfile)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to create CPU profile file")
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				logrus.Panic(err)
			}
		}(f)
		if err := pprof.StartCPUProfile(f); err != nil {
			logrus.WithError(err).Fatal("Failed to start CPU profile")
		}
		defer pprof.StopCPUProfile()
	}

	diBuilder, err := di.NewBuilder()

	// Config
	err = diBuilder.Add(di.Def{
		Name: "config",
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Parse(*flagConfigPath, "SUBAYAI_", config.DefaultConfig)
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("Config parsing failed")
	}

	// Database
	err = diBuilder.Add(di.Def{
		Name: "database",
		Build: func(ctn di.Container) (interface{}, error) {
			return database.New(ctn)
		},
		Close: func(obj interface{}) error {
			d := obj.(database.Database)
			logrus.Info("Shutting down database connection...")
			d.Close()
			return nil
		},
	})
	if err != nil && err.Error() == "unknown database driver" {
		logrus.WithError(err).Fatal("Database creation failed, unknown driver")
	} else if err != nil {
		logrus.WithError(err).Fatal("Database creation failed")
	}

	// Discord
	err = diBuilder.Add(di.Def{
		Name: "discord",
		Build: func(ctn di.Container) (interface{}, error) {
			return discord.New(ctn)
		},
		Close: func(obj interface{}) error {
			return obj.(*discord.Discord).Close()
		},
	})

	// Ken
	err = diBuilder.Add(di.Def{
		Name: "ken",
		Build: func(ctn di.Container) (interface{}, error) {
			return slashcommands.New(ctn)
		},
		Close: func(obj interface{}) error {
			return obj.(*ken.Ken).Unregister()
		},
	})

	ctn := diBuilder.Build()

	cfg := ctn.Get("config").(config.Config)
	level, err := logrus.ParseLevel(cfg.Logrus.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     cfg.Logrus.Colors,
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})

	logrus.Info("Starting Subayai ...")

	dc := ctn.Get("discord").(*discord.Discord)
	err = dc.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open Discord connection")
	}

	_ = ctn.Get("ken")

	block()

	err = ctn.DeleteWithSubContainers()
	if err != nil {
		return
	}

}

func block() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
