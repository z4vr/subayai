package leveling

import (
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/discord"
)

type Provider struct {
	dc       *discord.Discord
	db       database.Database
	guildMap GuildMap
}

func New(dc *discord.Discord, db database.Database) *Provider {
	return &Provider{
		dc:       dc,
		db:       db,
		guildMap: make(GuildMap),
	}
}

func (p *Provider) Open() {

	// populate guildMap
	session := p.dc.Session()
	guilds := session.State.Guilds

	for _, guild := range guilds {
		p.guildMap[guild.ID] = make(MemberMap)
	}

	// populate guildMap with members
	for _, guild := range guilds {
		for _, member := range guild.Members {
			p.guildMap[guild.ID][member.User.ID] = p.FetchFromDB(member.User.ID, guild.ID)
		}
	}

}

func (p *Provider) Close() {

}

func (p *Provider) FetchFromDB(guildID, userID string) *LevelData {
	levelData := &LevelData{
		UserID:  userID,
		GuildID: guildID,
	}

	var err error

	levelData.CurrentXP, err = p.db.GetUserCurrentXP(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": guildID,
			"userID":  userID,
		}).Error("Failed to fetch current XP from DB: ", err)
		return nil
	}
	levelData.TotalXP, err = p.db.GetUserTotalXP(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": guildID,
			"userID":  userID,
		}).Error("Failed to fetch total XP from DB: ", err)
		return nil
	}
	levelData.Level, err = p.db.GetUserLevel(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": guildID,
			"userID":  userID,
		}).Error("Failed to fetch level from DB: ", err)
		return nil
	}

	return levelData

}
