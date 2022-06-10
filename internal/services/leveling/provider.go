package leveling

import (
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
)

type Provider struct {
	db database.Database
}

func New(db database.Database) *Provider {
	return &Provider{
		db: db,
	}
}

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
