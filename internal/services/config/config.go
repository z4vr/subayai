package config

var DefaultConfig = Config{
	Discord: Discord{
		Token:      "",
		OwnerID:    "",
		GuildLimit: -1,
	},
	Logging: Logging{
		Level:  "info",
		Colors: true,
	},
	Database: Database{
		Type: "postgres",
		Postgres: DBCreds{
			Host: "localhost",
			Port: 5432,
		},
	},
	Waterlink: Waterlink{
		Host:     "",
		Password: "",
	},
}

type Discord struct {
	Token      string
	OwnerID    string
	GuildLimit int
}

type Logging struct {
	Level  string
	Colors bool
}

type Database struct {
	Type     string
	Postgres DBCreds
}

type DBCreds struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Waterlink struct {
	Host     string
	Password string
}

type Config struct {
	Discord   Discord
	Logging   Logging
	Database  Database
	Waterlink Waterlink
}
