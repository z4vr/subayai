package xp

import (
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/models"
	"github.com/z4vr/subayai/internal/services/database"
)

// GetUserXP returns the user's XP data.
func GetUserXP(userID, guildID string, db database.Database) (*models.UserXP, error) {

	level, err := db.GetUserLevel(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": guildID,
			"uid": userID,
		}).WithError(err).Error("Failed to get user level")
		return nil, err
	}
	currentXP, err := db.GetUserCurrentXP(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": guildID,
			"uid": userID,
		}).WithError(err).Error("Failed to get user current XP")
		return nil, err
	}
	totalXP, err := db.GetUserTotalXP(userID, guildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": guildID,
			"uid": userID,
		}).WithError(err).Error("Failed to get user total XP")
		return nil, err
	}

	xpData := models.New(userID, guildID, currentXP, totalXP, level)

	return xpData, nil
}

// GenerateUserXP generates the user's XP data with default values.
func GenerateUserXP(userID, guildID string, db database.Database) (*models.UserXP, error) {

	xpData := models.New(userID, guildID, 0, 0, 0)

	err := db.SetUserCurrentXP(userID, guildID, xpData.CurrentXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": guildID,
			"uid": userID,
		}).WithError(err).Error("Failed to set user current XP")
		return nil, err
	}
	err = db.SetUserTotalXP(userID, guildID, xpData.TotalXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": guildID,
			"uid": userID,
		}).WithError(err).Error("Failed to set user total XP")
		return nil, err
	}
	err = db.SetUserLevel(userID, guildID, xpData.Level)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": guildID,
			"uid": userID,
		}).WithError(err).Error("Failed to set user level")
		return nil, err
	}

	return xpData, nil

}

// UpdateUserXP updates the user's XP data.
func UpdateUserXP(xpData *models.UserXP, db database.Database) (err error) {

	err = db.SetUserCurrentXP(xpData.UserID, xpData.GuildID, xpData.CurrentXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": xpData.GuildID,
			"uid": xpData.UserID,
		}).WithError(err).Error("Failed to set user current XP")
		return err
	}
	err = db.SetUserTotalXP(xpData.UserID, xpData.GuildID, xpData.TotalXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": xpData.GuildID,
			"uid": xpData.UserID,
		}).WithError(err).Error("Failed to set user total XP")
		return err
	}
	err = db.SetUserLevel(xpData.UserID, xpData.GuildID, xpData.Level)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": xpData.GuildID,
			"uid": xpData.UserID,
		}).WithError(err).Error("Failed to set user level")
		return err
	}

	return nil
}
