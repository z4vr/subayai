package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

type ConfitaProvider struct {
	loc string
	*baseProvider
}

func NewConfitaProvider(loc string) *ConfitaProvider {
	return &ConfitaProvider{
		baseProvider: newBaseProvider(),
		loc:          loc,
	}
}

func (p *ConfitaProvider) Load() error {
	return confita.NewLoader(
		file.NewOptionalBackend(p.loc),
		env.NewBackend(),
		flags.NewBackend(),
	).Load(context.Background(), p.instance)
}
