package config

var DefaultConfig = Config{
	Bot: BotConfig{
		Token:      "",
		OwnerID:    "",
		GuildLimit: -1,
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
		Host:     "",
		Password: "",
	},
}

type BotConfig struct {
	Token      string `json:"token"`
	OwnerID    string `json:"ownerid"`
	GuildLimit int    `json:"guildlimit"`
}

type LogrusConfig struct {
	Level string `json:"level"`
	Color bool   `json:"color"`
}

type DatabaseConfig struct {
	Type     string         `json:"type"`
	Postgres PostgresConfig `json:"postgres"`
}

type PostgresConfig struct {
	Host     string `json:"postgres.host"`
	Port     int    `json:"postgres.port"`
	Database string `json:"postgres.database"`
	Username string `json:"postgres.username"`
	Password string `json:"postgres.password"`
}

type LavalinkConfig struct {
	Host     string `json:"host"`
	Password string `json:"password"`
}

type Config struct {
	Bot      BotConfig      `json:"bot"`
	Logrus   LogrusConfig   `json:"logrus"`
	Database DatabaseConfig `json:"database"`
	Lavalink LavalinkConfig `json:"lavalink"`
}
