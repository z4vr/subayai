package inits

import (
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/database/postgres"
	"github.com/z4vr/subayai/internal/util/static"
)

func InitDatabase(ctn di.Container) database.Database {
	var db database.Database
	var err error

	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	switch cfg.Config().Database.Type {
	case "postgres":
		db = new(postgres.PGMiddleware)
		err = db.Connect(cfg.Config().Database.Postgres)
	default:
		logrus.Fatalf("Unknown database type: %s", cfg.Config().Database.Type)
	}
	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting to database")
	}

	return db
}
