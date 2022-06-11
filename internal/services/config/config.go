package config

var DefaultConfig = Config{
	Discord: Discord{
		Token:      "",
		OwnerID:    "",
		GuildLimit: -1,
	},
	Logrus: Logrus{
		Level:  "info",
		Colors: true,
	},
	Database: Database{
		Type: "postgres",
		DBCreds: DBCreds{
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

type Logrus struct {
	Level  string
	Colors bool
}

type Database struct {
	Type    string
	DBCreds DBCreds
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
	Logrus    Logrus
	Database  Database
	Waterlink Waterlink
}
