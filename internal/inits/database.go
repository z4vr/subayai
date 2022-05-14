package inits

import (
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/config"
	"github.com/z4vr/subayai/internal/database"
	"github.com/z4vr/subayai/internal/database/fio"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/pkg/fileio"
)

func InitDatabase(ctn di.Container) database.Database {
	var db database.Database
	var err error
	var dbTables = []string{
		"guilds",
		"xp",
		"users",
	}

	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)

	db = new(fio.MiddleWare)
	err = db.Connect(cfg.Instance().FIO)
	fProvider := db.RawProvider().(*fileio.Provider)
	if !fProvider.CheckFolder(fProvider.FIOPath) {
		err = fProvider.GenerateFolder(fProvider.FIOPath)
	}
	err = db.CreateTables(dbTables)

	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting to database")
	}

	return db
}
