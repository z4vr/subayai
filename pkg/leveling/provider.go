package leveling

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/errorarray"
)

type LevelProvider struct {
	session *discordgo.Session
	db      database.Database
	guilds  GuildData
}

func (l *LevelProvider) Close() error {
	// save the guild data
	err := l.SaveGuildData(l.session.State.Guilds)
	if err != nil {
		return err
	}

	// delete the guild data
	l.guilds = make(GuildData)

	return err
}

func New(session *discordgo.Session, db database.Database) *LevelProvider {
	return &LevelProvider{
		session: session,
		db:      db,
		guilds:  make(GuildData),
	}
}

// Get returns the data for a given user in a guild
func (l *LevelProvider) Get(userID, guildID string) *Data {
	return l.guilds[guildID][userID]
}

// Set sets the data for a given user in a guild
func (l *LevelProvider) Set(userID, guildID string, d *Data) {
	l.guilds[guildID][userID] = d
}

// PopulateGuildData populates the guild data with the data from the database
func (l *LevelProvider) PopulateGuildData(guilds []*discordgo.Guild) error {
	errArray := errorarray.New()
	for _, guild := range guilds {
		guildData := l.guilds[guild.ID]
		for _, member := range guild.Members {
			d := l.FetchFromDB(member.User.ID, guild.ID)
			if d != nil {
				guildData[member.User.ID] = d
			}
		}
		l.guilds[guild.ID] = guildData
	}
	return errArray.Nillify()
}

// FetchFromDB returns the leveling data for the given user
func (l *LevelProvider) FetchFromDB(userID, guildID string) *Data {

	currentXP, err := l.db.GetUserCurrentXP(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user current XP")
		return nil
	}
	totalXP, err := l.db.GetUserTotalXP(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user total XP")
		return nil
	}
	level, err := l.db.GetUserLevel(userID, guildID)
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
	guildData := l.guilds[d.GuildID]
	guildData[d.UserID] = d
	l.guilds[d.GuildID] = guildData

	return l.SaveLevelData(*d)
}

// Remove removes the given Data from the guild data
func (l *LevelProvider) Remove(d Data) {
	guildData := l.guilds[d.GuildID]
	delete(guildData, d.UserID)
	l.guilds[d.GuildID] = guildData
}

func (l *LevelProvider) GetGuildData(guildID string) UserData {
	return l.guilds[guildID]
}

func (l *LevelProvider) GetUserData(userID, guildID string) *Data {
	return l.guilds[guildID][userID]
}

func (l *LevelProvider) SaveGuildData(guilds []*discordgo.Guild) error {
	errArray := errorarray.New()
	for _, guild := range guilds {
		guildData := l.guilds[guild.ID]
		for _, member := range guild.Members {
			d := guildData[member.User.ID]
			errArray.Append(l.SaveLevelData(*d))
		}
	}
	return errArray.Nillify()
}

func (l *LevelProvider) SaveLevelData(d Data) error {
	err := l.db.SetUserCurrentXP(d.UserID, d.GuildID, d.CurrentXP)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user current XP")
		return err
	}
	err = l.db.SetUserTotalXP(d.UserID, d.GuildID, d.TotalXP)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user total XP")
		return err
	}
	err = l.db.SetUserLevel(d.UserID, d.GuildID, d.Level)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user leveling")
		return err
	}

	return nil

}
