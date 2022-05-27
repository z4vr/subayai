package level

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/util/static"
)

type LevelProvider struct {
	session *discordgo.Session
	db      database.Database
	guilds  GuildData
}

func NewLevelProvider(ctn di.Container) *LevelProvider {
	s := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	db := ctn.Get(static.DiDatabase).(database.Database)

	provider := &LevelProvider{
		session: s,
		db:      db,
	}

	// iterate over the current guilds
	// and load the data from the db into the map
	guilds := s.State.Guilds
	if len(guilds) == 0 {
		// assign an empty map to
		// the provider
		provider.guilds = make(GuildData)

		return provider

	}

	guildData := make(GuildData, len(guilds))

	// iterate over the guilds
	for _, guild := range guilds {
		// create a new map for the guild
		guildData[guild.ID] = make(UserData, len(guild.Members))

		// iterate over the members
		for _, member := range guild.Members {
			levelData := provider.GetLevelData(member.User.ID, guild.ID)
			if levelData == nil {
				continue
			}
			// add the data to the map
			guildData[guild.ID][member.User.ID] = *levelData
		}
	}

	provider.guilds = guildData

	return provider

}

// GetLevelData returns the level data for the given user
func (p *LevelProvider) GetLevelData(userID, guildID string) *Data {

	currentXP, err := p.db.GetUserCurrentXP(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user current XP")
		return nil
	}
	totalXP, err := p.db.GetUserTotalXP(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user total XP")
		return nil
	}
	level, err := p.db.GetUserLevel(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user level")
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
func (p *LevelProvider) Add(d Data) error {
	guildData := p.guilds[d.GuildID]
	guildData[d.UserID] = d
	p.guilds[d.GuildID] = guildData

	return p.SaveLevelData(d)
}

// Remove removes the given Data from the guild data
func (p *LevelProvider) Remove(d Data) {
	guildData := p.guilds[d.GuildID]
	delete(guildData, d.UserID)
	p.guilds[d.GuildID] = guildData
}

func (p *LevelProvider) GetGuildData(guildID string) UserData {
	return p.guilds[guildID]
}

func (p *LevelProvider) GetUserData(userID, guildID string) Data {
	return p.guilds[guildID][userID]
}

func (p *LevelProvider) SaveLevelData(data Data) error {
	err := p.db.SetUserCurrentXP(data.UserID, data.GuildID, data.CurrentXP)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user current XP")
		return err
	}
	err = p.db.SetUserTotalXP(data.UserID, data.GuildID, data.TotalXP)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user total XP")
		return err
	}
	err = p.db.SetUserLevel(data.UserID, data.GuildID, data.Level)
	if err != nil {
		logrus.WithError(err).Error("Failed to save user level")
		return err
	}

	return nil

}
