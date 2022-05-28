package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
)

var (
	flagConfigPath = flag.String("config", "config.yaml", "Path to config file")
	flagCPUProfile = flag.String("cpuprofile", "", "Path to write CPU profile")
)

func main() {
	flag.Parse()

	cfg := config.NewPaerser(*flagConfigPath)
	if err := cfg.Parse(); err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}
	loglevel, err := logrus.ParseLevel(cfg.Config().Logrus.Level)
	if err != nil {
		logrus.WithError(err).Warn("Failed to parse logrus leveling, using default")
		loglevel = logrus.InfoLevel
	}
	logrus.SetLevel(loglevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     cfg.Config().Logrus.Color,
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})

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

	block()
}

func block() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
