package database

import "github.com/z4vr/subayai/pkg/database/postgres"

type Config struct {
	Type     string
	Postgres postgres.Config
}
