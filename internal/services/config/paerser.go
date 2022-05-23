package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/file"
	"github.com/z4vr/subayai/internal/models"
)

type Paerser struct {
	cfg        *models.Config
	configFile string
}

func NewPaerser(configFile string) *Paerser {
	return &Paerser{
		configFile: configFile,
	}
}

func (p *Paerser) Config() *models.Config {
	return p.cfg
}

func (p *Paerser) Parse() (err error) {
	cfg := models.DefaultConfig

	if err = file.Decode(p.configFile, &cfg); err != nil && !os.IsNotExist(err) {
		return
	}

	godotenv.Load()
	if err = env.Decode(os.Environ(), "SP_", &cfg); err != nil {
		return
	}

	p.cfg = &cfg

	return
}
