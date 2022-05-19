package config

import (
	"os/user"

	"github.com/z4vr/subayai/internal/models"
)

var currentUser, _ = user.Current()

type Provider interface {
	Load() error
	Instance() *models.Config
}

type baseProvider struct {
	instance *models.Config
}

func newBaseProvider() *baseProvider {
	return &baseProvider{
		instance: models.DefaultConfig,
	}
}

func (p *baseProvider) Instance() *models.Config {
	return p.instance
}
