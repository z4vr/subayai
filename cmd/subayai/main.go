package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/discord"
	"github.com/z4vr/subayai/internal/services/leveling"

	"github.com/sirupsen/logrus"
)

var (
	flagConfigPath = flag.String("config", "config.toml", "Path to config file")
	flagCPUProfile = flag.String("cpuprofile", "", "Path to write CPU profile")
)

func main() {

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	// Config
	cfg, err := config.Parse(*flagConfigPath, "SUBAYAI_", config.DefaultConfig)
	if err != nil {
		logrus.WithError(err).Fatal("Config parsing failed")
	}

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

	// Database
	db := database.New(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Database initialization failed")
	}
	defer func() {
		logrus.Info("Shutting down database connection ...")
		db.Close()
	}()
	logrus.WithField("type", cfg.Database.Type).Info("Database initialized")

	// Discord & Leveling
	lp := leveling.New(db)
	if lp == nil {
		logrus.Fatal("Leveling initialization failed")
	}
	dc, err := discord.New(cfg, db, lp)
	if err != nil {
		logrus.WithError(err).Fatal("Discord initialization failed")
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

	block()
}

func block() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
