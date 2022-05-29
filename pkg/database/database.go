package database

import (
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database/postgres"
)

// Database is the interface for a database driver.
type Database interface {

	// GeneralL
	Connect(credentials ...interface{}) error

	// guilds
	GetGuildBotMessageChannelID(guildID string) (channelID string, err error)
	SetGuildBotMessageChannelID(guildID, channelID string) error

	GetGuildLevelUpMessage(guildID string) (message string, err error)
	SetGuildLevelUpMessage(guildID, message string) (err error)

	GetGuildAFKChannelID(guildID string) (channelID string, err error)
	SetGuildAFKChannelID(guildID, channelID string) (err error)

	GetGuildAutoroleIDs(guildID string) (roleIDs []string, err error)
	SetGuildAutoroleIDs(guildID string, roleIDs []string) (err error)

	GetUserLevel(userID, guildID string) (level int, err error)
	SetUserLevel(userID, guildID string, level int) (err error)

	GetUserCurrentXP(userID, guildID string) (xp int, err error)
	SetUserCurrentXP(userID, guildID string, xp int) (err error)

	GetUserTotalXP(userID, guildID string) (xp int, err error)
	SetUserTotalXP(userID, guildID string, xp int) (err error)

	GetLastMessageTimestamp(userID, guildID string) (timestamp int64, err error)
	SetLastMessageTimestamp(userID, guildID string, timestamp int64) (err error)

	GetLastVoiceSessionTimestamp(userID, guildID string) (timestamp int64, err error)
	SetLastVoiceSessionTimestamp(userID, guildID string, timestamp int64) (err error)

	GetLastMessageID(userID, guildID string) (id string, err error)
	SetLastMessageID(userID, guildID string, id string) (err error)

	GetLastVoiceSessionID(userID, guildID string) (id string, err error)
	SetLastVoiceSessionID(userID, guildID string, id string) (err error)

	Close()
}

func New(cfg Config) Database {
	var db Database
	var err error

	switch cfg.Type {
	case "postgres":
		db = new(postgres.PGMiddleware)
		err = db.Connect(cfg.Postgres)
	default:
		logrus.Fatalf("Unknown database type: %s", cfg.Type)
	}
	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting to database")
	}

	return db
}