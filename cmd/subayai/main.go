package main

import (
	"flag"
	"io"
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

	logFile, err := os.OpenFile("./subayai.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     false,
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})
	logrus.SetOutput(mw)

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
	lp := leveling.New(db)
	if lp == nil {
		logrus.Fatal("Leveling initialization failed")
	}
	dc, err := discord.New(cfg.Discord, db, lp)
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
