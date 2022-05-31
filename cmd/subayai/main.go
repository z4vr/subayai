package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/z4vr/subayai/pkg/config"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/discord"
	"github.com/z4vr/subayai/pkg/leveling"

	"github.com/sirupsen/logrus"
)

var (
	flagConfigPath = flag.String("config", "config.yaml", "Path to config file")
	flagCPUProfile = flag.String("cpuprofile", "", "Path to write CPU profile")
)

func main() {
	flag.Parse()

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})

	// Config
	cfg, err := config.Parse(*flagConfigPath, "SUBAYAI_", config.DefaultConfig)
	if err != nil {
		logrus.WithError(err).Fatal("Config parsing failed")
	}

	if *flagCPUProfile != "" {
		f, err := os.Create(*flagCPUProfile)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to create CPU profile file")
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			logrus.WithError(err).Fatal("Failed to start CPU profile")
		}
		defer pprof.StopCPUProfile()
	}

	// Database
	db := database.New(cfg.Database)
	if err != nil {
		logrus.WithError(err).Fatal("Database initialization failed")
	}
	defer func() {
		logrus.Info("Shutting down database connection ...")
		db.Close()
	}()
	logrus.WithField("typ", cfg.Database.Type).Info("Database initialized")

	// Discord & Leveling
	dc, err := discord.New(cfg.Discord, db)
	if err != nil {
		logrus.WithError(err).Fatal("Discord initialization failed")
	}
	lp := leveling.New(dc, db)
	if lp == nil {
		logrus.Fatal("Leveling initialization failed")
	}
	err = dc.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Discord connection failed")
	}
	logrus.Info("Discord connection initialized")
	defer func() {
		logrus.Info("Shutting down Discord connection ...")
		dc.Close()
	}()
	lp.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Leveling connection failed")
	}
	logrus.Info("Leveling map initialized")
	defer func() {
		logrus.Info("Shutting down leveling map ...")
		err := lp.Close()
		if err != nil {
			logrus.WithError(err).Fatal("Leveling map shutdown failed -> check database entrys for corruption")
		}
	}()

	block()
}

func block() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
