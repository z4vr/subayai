package leveling

import (
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/discord"
	"github.com/z4vr/subayai/pkg/errorarray"
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

// Open and Close

func (p *Provider) Open() {

	// create an empty guildMap
	p.guildMap = make(GuildMap)

	// create a map for each guild#
	// we don't want to create a map for every user
	// since this might take too long to load
	// and we don't want to load all users at once
	for _, guild := range p.dc.Session().State.Guilds {
		p.guildMap[guild.ID] = make(MemberMap)
	}

	// add listeners
	p.dc.Session().AddHandler(p.MessageCreate)

}

func (p *Provider) Close() error {

	errorArray := errorarray.New()
	// save guildMap to DB
	for _, memberMap := range p.guildMap {
		for _, levelData := range memberMap {
			err := p.SaveToDB(levelData)
			if err != nil {
				errorArray.Append(err)
				continue
			}
		}
	}

	if errorArray.Len() > 0 {
		logrus.Error("Failed to save level data to DB: ", errorArray.Errors())
		return errorArray
	}

	return nil

}

// Getter and Setters

func (p *Provider) Get(userID, guildID string) (*LevelData, error) {

	if p.guildMap[guildID][userID] == nil {
		return p.FetchFromDB(userID, guildID)
	}

	return p.guildMap[guildID][userID], nil

}

func (p *Provider) Set(userID, guildID string, levelData *LevelData) {
	p.guildMap[guildID][userID] = levelData
}

// DB functions

func (p *Provider) SaveToDB(levelData *LevelData) error {

	err := p.db.SetUserCurrentXP(levelData.UserID, levelData.GuildID, levelData.CurrentXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).Error("Failed to save current XP to DB: ", err)
		return err
	}
	err = p.db.SetUserTotalXP(levelData.UserID, levelData.GuildID, levelData.TotalXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).Error("Failed to save total XP to DB: ", err)
		return err
	}
	err = p.db.SetUserLevel(levelData.UserID, levelData.GuildID, levelData.Level)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).Error("Failed to save level to DB: ", err)
		return err
	}

	return err

}

func (p *Provider) FetchFromDB(userID, guildID string) (*LevelData, error) {
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
		return nil, err
	}
	levelData.TotalXP, err = p.db.GetUserTotalXP(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": guildID,
			"userID":  userID,
		}).Error("Failed to fetch total XP from DB: ", err)
		return nil, err
	}
	levelData.Level, err = p.db.GetUserLevel(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": guildID,
			"userID":  userID,
		}).Error("Failed to fetch level from DB: ", err)
		return nil, err
	}

	return levelData, nil

}
