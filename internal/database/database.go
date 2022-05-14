package database

import "github.com/z4vr/subayai/internal/models"

// Database is the interface for a database driver.
type Database interface {

	// GENERAL
	Connect(credentials ...interface{}) error
	RawProvider() interface{}
	CreateTables([]string) error

	// GUILD
	GetGuildConfig(guildID string) (models.GuildConfig, error)
	SetGuildConfig(guildID string, config models.GuildConfig) error

	GetGuildAutoroleIDs(guildID string) ([]string, error)
	SetGuildAutoroleIDs(guildID string, roleIDs []string) error

	GetGuildAutoDelete(guildID string) (bool, error)
	SetGuildAutoDelete(guildID string, autoDelete bool) error

	GuildEntryExists(guildID string) (bool, error)
	CreateGuildEntry(guildID string) error
	DeleteGuildEntry(guildID string) error

	// USER
	GetUserConfig(userID string) (models.UserConfig, error)
	SetUserConfig(userID string, config models.UserConfig) error

	UserEntryExists(userID string) (bool, error)
	CreateUserEntry(userID string) error
	DeleteUserEntry(userID string) error

	// XP
	GetUserXPEntry(userID, guildID string) (models.UserXPEntry, error)
	SetUserXPEntry(userID, guildID string, entry models.UserXPEntry) error

	GetUserLevel(userID, guildID string) (int, error)
	SetUserLevel(userID, guildID string, level int) error

	GetUserCurrentXP(userID, guildID string) (int, error)
	SetUserCurrentXP(userID, guildID string, currentXP int) error

	GetUserTotalXP(userID, guildID string) (int, error)
	SetUserTotalXP(userID, guildID string, totalXP int) error

	GetUserLastMessageTimestamp(userID, guildID string) (int64, error)
	SetUserLastMessageTimestamp(userID, guildID string, lastMessageTimestamp int64) error

	GetUserLastSessionID(userID, guildID string) (string, error)
	SetUserLastSessionID(userID, guildID string, lastSessionID string) error

	GetUserLastSessionTimestamp(userID, guildID string) (int64, error)
	SetUserLastSessionTimestamp(userID, guildID string, lastSessionTimestamp int64) error

	UserXPEntryExists(userID, guildID string) (bool, error)
	CreateUserXPEntry(userID, guildID string) error
	DeleteUserXPEntry(userID, guildID string) error
}
