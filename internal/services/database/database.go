package database

import "errors"

var ErrValueNotFound = errors.New("value not found in database")

// Database is the interface for a database driver.
type Database interface {

	// GeneralL
	Connect(credentials ...interface{}) error

	// Guilds
	GetGuildBotMessageChannelID(guildID string) (channelID string, err error)
	SetGuildBotMessageChannelID(guildID, channelID string) error

	GetGuildLevelUpMessage(guildID string) (message string, err error)
	SetGuildLevelUpMessage(guildID, message string) error

	GetGuildAFKChannelID(guildID string) (channelID string, err error)
	SetGuildAFKChannelID(guildID, channelID string) error

	GetGuildAutoroleIDs(guildID string) (roleIDs []string, err error)
	SetGuildAutoroleIDs(guildID string, roleIDs []string) error

	Close()
}
