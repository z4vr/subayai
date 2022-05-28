package leveling

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/pkg/errorarray"
)

type LevelProvider struct {
	Session *discordgo.Session
	Db      database.Database
	Guilds  GuildData
}

func (l *LevelProvider) Close() error {
	// save the guild data
	err := l.SaveGuildData(l.Session.State.Guilds)
	if err != nil {
		return err
	}

	// delete the guild data
	l.Guilds = make(GuildData)

	return err
}

// Get returns the data for a given user in a guild
func (l *LevelProvider) Get(userID, guildID string) *Data {
	return l.Guilds[guildID][userID]
}

// Set sets the data for a given user in a guild
func (l *LevelProvider) Set(userID, guildID string, d *Data) {
	l.Guilds[guildID][userID] = d
}

// FetchFromDB returns the leveling data for the given user
func (l *LevelProvider) FetchFromDB(userID, guildID string) *Data {

	currentXP, err := l.Db.GetUserCurrentXP(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user current XP")
		return nil
	}
	totalXP, err := l.Db.GetUserTotalXP(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user total XP")
		return nil
	}
	level, err := l.Db.GetUserLevel(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user leveling")
		return nil
	}

	return &Data{
		UserID:    userID,
		GuildID:   guildID,
		CurrentXP: currentXP,
		TotalXP:   totalXP,
		Level:     level,
	}

}

// Add appends the given Data to the guild data
func (l *LevelProvider) Add(d *Data) error {
	guildData := l.Guilds[d.GuildID]
	guildData[d.UserID] = d
	l.Guilds[d.GuildID] = guildData

	return l.SaveLevelData(*d)
}

// Remove removes the given Data from the guild data
func (l *LevelProvider) Remove(d Data) {
	guildData := l.Guilds[d.GuildID]
	delete(guildData, d.UserID)
	l.Guilds[d.GuildID] = guildData
}

func (l *LevelProvider) GetGuildData(guildID string) UserData {
	return l.Guilds[guildID]
}

func (l *LevelProvider) GetUserData(userID, guildID string) *Data {
	return l.Guilds[guildID][userID]
}

func (l *LevelProvider) SaveGuildData(guilds []*discordgo.Guild) error {
	errArray := errorarray.New()
	for _, guild := range guilds {
		guildData := l.Guilds[guild.ID]
		for _, member := range guild.Members {
			d := guildData[member.User.ID]
			errArray.Append(l.SaveLevelData(*d))
		}
	}
	return errArray.Nillify()
}

func (l *LevelProvider) SaveLevelData(d Data) error {
	err := l.Db.SetUserCurrentXP(d.UserID, d.GuildID, d.CurrentXP)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user current XP")
		return err
	}
	err = l.Db.SetUserTotalXP(d.UserID, d.GuildID, d.TotalXP)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user total XP")
		return err
	}
	err = l.Db.SetUserLevel(d.UserID, d.GuildID, d.Level)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user leveling")
		return err
	}

	return nil

}
