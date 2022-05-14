package config

import (
	"fmt"
	"os/user"

	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/models"
)

var currentUser, _ = user.Current()

var defaultConfig = models.Config{
	Discord: models.Discord{},
	Log: models.Log{
		Level:  logrus.InfoLevel,
		Colors: true,
	},
	Lavalink: models.Lavalink{
		Authorization: "",
		Host:          "localhost",
		Port:          80,
	},
	FIO: models.FIO{
		FIOPath: fmt.Sprintf("%s/.config/subayai/fio", currentUser.HomeDir),
	},
}

type Provider interface {
	Load() error
	Instance() *models.Config
}

type baseProvider struct {
	instance *models.Config
}

func newBaseProvider() *baseProvider {
	return &baseProvider{
		instance: &defaultConfig,
	}
}

func (p *baseProvider) Instance() *models.Config {
	return p.instance
}
