package database

import "errors"

var ErrValueNotFound = errors.New("value not found in database")

// Database is the interface for a database driver.
type Database interface {

	// GeneralL
	Connect(credentials ...interface{}) error
	Close()

	// Guilds
	GetGuildBotMessageChannelID(guildID string) (channelID string, err error)
	SetGuildBotMessageChannelID(guildID, channelID string) error

	GetGUildLevelUpMessage(guildID string) (message string, err error)
	SetGuildLevelUpMessage(guildID, message string) error

	GetGuildAFKChannelID(guildID string) (channelID string, err error)
	SetGuildAFKChannelID(guildID, channelID string) error

	GetGuildAutoRolesID(guildID string) (roleID []string, err error)
	SetGuildAutoRolesID(guildID, roleID []string) error
}
