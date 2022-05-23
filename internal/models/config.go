package models

var DefaultConfig = Config{
	Discord: DiscordConfig{
		Token:   "",
		OwnerID: "",
	},
	Logrus: LogrusConfig{
		Level: "info",
		Color: true,
	},
	Database: DatabaseConfig{
		Type: "postgres",
		Postgres: PostgresConfig{
			Host: "localhost",
			Port: 5432,
		},
	},
	Lavalink: LavalinkConfig{
		Authorization: "",
		Host:          "",
	},
}

type DiscordConfig struct {
	Token   string `config:"token"`
	OwnerID string `config:"ownerid"`
}

type LogrusConfig struct {
	Level string `config:"level"`
	Color bool   `config:"color"`
}

type DatabaseConfig struct {
	Type     string         `config:"type"`
	Postgres PostgresConfig `config:"postgres"`
}

type PostgresConfig struct {
	Host     string `config:"postgres.host"`
	Port     int    `config:"postgres.port"`
	Database string `config:"postgres.database"`
	Username string `config:"postgres.username"`
	Password string `config:"postgres.password"`
}

type LavalinkConfig struct {
	Authorization string `config:"lavalink.authorization"`
	Host          string `config:"lavalink.host"`
}

type Config struct {
	Discord  DiscordConfig  `config:"discord"`
	Logrus   LogrusConfig   `config:"logrus"`
	Database DatabaseConfig `config:"database"`
	Lavalink LavalinkConfig `config:"lavalink"`
}
