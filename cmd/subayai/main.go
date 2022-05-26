package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/z4vr/subayai/internal/util/static"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/inits"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database"
)

var (
	flagConfigPath = flag.String("config", "./config.yaml", "Path to config file")
	flagCPUProfile = flag.String("cpuprofile", "", "Path to write CPU profile")
)

func main() {
	flag.Parse()

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

	diBuilder, err := di.NewBuilder()
	if err != nil {
		logrus.WithError(err).Panic("Error initializing DI")
	}

	// Register dependencies
	// register config
	err = diBuilder.Add(di.Def{
		Name: static.DiConfigProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewPaerser(*flagConfigPath), nil
		},
	})
	if err != nil {
		logrus.WithError(err).Panic("Error initializing Config Provider")
		return
	}

	//Register database
	err = diBuilder.Add(di.Def{
		Name: static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDatabase(ctn), nil
		},
	})
	if err != nil {
		logrus.WithError(err).Panic("Error initializing Database")
		return
	}

	// register discord session
	err = diBuilder.Add(di.Def{
		Name: static.DiDiscordSession,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.NewDiscordSession(ctn)
		},
		Close: func(obj interface{}) error {
			logrus.Info("Shutting down Bot connection...")
			return obj.(*discordgo.Session).Close()
		},
	})
	if err != nil {
		logrus.WithError(err).Panic("Error initializing Discord Session")
		return
	}

	// Building object map
	ctn := diBuilder.Build()
	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	if err := cfg.Parse(); err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}
	level, err := logrus.ParseLevel(cfg.Config().Logrus.Level)
	if err != nil {
		logrus.WithError(err).Warn("Failed to parse logrus level, using default")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     cfg.Config().Logrus.Color,
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})

	session := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	if err := session.Open(); err != nil {
		logrus.WithError(err).Fatal("Failed connecting to discord")
	}

	// get the database object to initialize the connection
	_ = ctn.Get(static.DiDatabase).(database.Database)

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
