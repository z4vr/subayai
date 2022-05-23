package config

import "github.com/z4vr/subayai/internal/models"

type Provider interface {
	Config() *models.Config
	Parse() error
}
