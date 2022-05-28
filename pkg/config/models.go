package config

import (
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/database/postgres"
	"github.com/z4vr/subayai/pkg/discord"
	"github.com/z4vr/subayai/pkg/waterlink"
)

var DefaultConfig = Config{
	Discord: discord.Config{
		Token:      "",
		OwnerID:    "",
		GuildLimit: -1,
	},
	Database: database.Config{
		Type: "postgres",
		Postgres: postgres.Config{
			Host: "localhost",
			Port: 5432,
		},
	},
	Lavalink: waterlink.Config{
		Host:     "",
		Password: "",
	},
}

type Config struct {
	Discord  discord.Config
	Database database.Config
	Lavalink waterlink.Config
}
