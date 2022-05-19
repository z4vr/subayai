package models

import (
	"github.com/sirupsen/logrus"
)

var DefaultConfig = &Config{
	Discord: Discord{
		Token:   "",
		OwnerId: "",
	},
	Log: Log{
		Level:  logrus.InfoLevel,
		Colors: true,
	},
	Database: DatabaseType{
		Type:     "postgres",
		Postgres: DatabaseCreds{},
	},
	Lavalink: Lavalink{
		Authorization: "",
		Host:          "",
	},
}

type Discord struct {
	Token   string `config:"token,required"`
	OwnerId string `config:"ownerid"`
}

type Log struct {
	Level  logrus.Level `config:"level"`
	Colors bool         `config:"colors"`
}

type DatabaseType struct {
	Type     string        `config:"type"`
	Postgres DatabaseCreds `config:"postgres"`
}

type DatabaseCreds struct {
	Host     string `config:"host"`
	User     string `config:"user"`
	Password string `config:"password"`
	Database string `config:"database"`
}

type Lavalink struct {
	Authorization string `config:"authorization"`
	Host          string `config:"host"`
}

type Config struct {
	Discord  Discord      `config:"discord"`
	Log      Log          `config:"log"`
	Database DatabaseType `json:"database"`
	Lavalink Lavalink     `config:"lavalink"`
}
