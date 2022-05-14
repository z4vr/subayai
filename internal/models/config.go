package models

import (
	"github.com/sirupsen/logrus"
)

type Discord struct {
	Token   string `config:"token,required"`
	OwnerId string `config:"ownerid"`
}

type Log struct {
	Level  logrus.Level `config:"level"`
	Colors bool         `config:"colors"`
}

type Lavalink struct {
	Authorization string `config:"authorization"`
	Host          string `config:"host"`
	Port          int    `config:"port"`
}

type FIO struct {
	FIOPath string `config:"fiopath"`
}

type Config struct {
	Discord  Discord  `config:"discord"`
	Log      Log      `config:"log"`
	Lavalink Lavalink `config:"lavalink"`
	FIO      FIO      `config:"fileio"`
}
