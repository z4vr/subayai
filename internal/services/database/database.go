package database

import (
	"errors"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database/postgres"
)

// Database is the interface for a database driver.
type Database interface {

	// General

	Connect(credentials ...interface{}) error

	// GUILDS

	GetGuildBotMessageChannelID(guildID string) (channelID string, err error)
	SetGuildBotMessageChannelID(guildID, channelID string) error

	GetGuildLevelUpMessage(guildID string) (message string, err error)
	SetGuildLevelUpMessage(guildID, message string) (err error)

	GetGuildAFKChannelID(guildID string) (channelID string, err error)
	SetGuildAFKChannelID(guildID, channelID string) (err error)

	GetGuildAutoroleIDs(guildID string) (roleIDs []string, err error)
	SetGuildAutoroleIDs(guildID string, roleIDs []string) (err error)

	// LEVELING

	GetUserLevel(userID, guildID string) (level int, err error)
	SetUserLevel(userID, guildID string, level int) (err error)

	GetUserCurrentXP(userID, guildID string) (xp int, err error)
	SetUserCurrentXP(userID, guildID string, xp int) (err error)

	GetUserTotalXP(userID, guildID string) (xp int, err error)
	SetUserTotalXP(userID, guildID string, xp int) (err error)

	// TIMESTAMPS

	GetLastMessageTimestamp(userID, guildID string) (timestamp int64, err error)
	SetLastMessageTimestamp(userID, guildID string, timestamp int64) (err error)

	GetLastVoiceSessionTimestamp(userID, guildID string) (timestamp int64, err error)
	SetLastVoiceSessionTimestamp(userID, guildID string, timestamp int64) (err error)

	// IDS

	GetLastMessageID(userID, guildID string) (id string, err error)
	SetLastMessageID(userID, guildID string, id string) (err error)

	GetLastVoiceSessionID(userID, guildID string) (id string, err error)
	SetLastVoiceSessionID(userID, guildID string, id string) (err error)

	GetRewardRoleIDByLevel(guildID string, level int) (id string, err error)
	SetRewardRoleIDByLevel(guildID string, level int, id string) (err error)

	Close()
}

func New(ctn di.Container) (Database, error) {
	var db Database
	var err error

	cfg := ctn.Get("config").(config.Config)

	switch cfg.Database.Type {
	case "postgres":
		db = new(postgres.PGMiddleware)
		err = db.Connect(cfg.Database.DBCreds)
	default:
		logrus.Fatalf("Unknown database type: %s", cfg.Database.Type)
		return nil, errors.New("unknown database type")
	}
	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting to database")
		return nil, err
	}

	return db, nil
}
