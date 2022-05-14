package fio

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/database"
	"github.com/z4vr/subayai/internal/models"
	"github.com/z4vr/subayai/pkg/fileio"
)

type MiddleWare struct {
	fio *fileio.Provider
}

var (
	// implemenation of the database interface
	_ database.Database = (*MiddleWare)(nil)
)

// Connect connects to the database. Here it just initializes the FIO provider.
func (f *MiddleWare) Connect(credentials ...interface{}) error {
	creds := credentials[0].(models.FIO)
	f.fio = fileio.NewFIOProvider(creds.FIOPath)
	return nil
}

// RawProvider returns the raw FIO provider.
func (f *MiddleWare) RawProvider() interface{} {
	return f.fio
}

// CreateTables creates the tables.
func (f *MiddleWare) CreateTables(tables []string) error {
	for _, table := range tables {
		ok := f.fio.CheckFolder(fmt.Sprintf("%s/%s", f.fio.FIOPath, table))
		if !ok {
			err := f.fio.GenerateFolder(fmt.Sprintf("%s/%s", f.fio.FIOPath, table))
			if err != nil {
				logrus.WithError(err).Error("Failed to create table: " + table)
				return err
			}
		}
	}
	return nil
}

// GetGuildConfig returns the guild config for a guild.
func (f *MiddleWare) GetGuildConfig(guildID string) (models.GuildConfig, error) {
	guildPath := fmt.Sprintf("%s/guilds/%s.json", f.fio.FIOPath, guildID)

	gMap, err := f.fio.Parse(guildPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse guild config")
		return models.GuildConfig{}, err
	}

	return models.GuildConfig{
		AutoDelete:  gMap["auto_delete"].(bool),
		AutoRoleIDs: gMap["auto_role_ids"].(string),
		GuildID:     gMap["guild_id"].(string),
	}, err
}

// SetGuildConfig sets the guild config for a guild.
func (f *MiddleWare) SetGuildConfig(guildID string, guildConfig models.GuildConfig) error {
	guildPath := fmt.Sprintf("%s/guilds/%s.json", f.fio.FIOPath, guildID)

	gInterface := map[string]interface{}{
		"auto_delete":   guildConfig.AutoDelete,
		"auto_role_ids": guildConfig.AutoRoleIDs,
		"guild_id":      guildConfig.GuildID,
	}

	return f.fio.Save(guildPath, gInterface)
}

// GetGuildAutoroleIDs returns the autorole IDs for a guild.
func (f *MiddleWare) GetGuildAutoroleIDs(guildID string) ([]string, error) {

	guildConfig, err := f.GetGuildConfig(guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get guild config")
		return []string{}, err
	}

	return strings.Split(guildConfig.AutoRoleIDs, ";"), nil

}

// SetGuildAutoroleIDs sets the autorole IDs for a guild.
func (f *MiddleWare) SetGuildAutoroleIDs(guildID string, roleIDs []string) error {

	guildConfig, err := f.GetGuildConfig(guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get guild config")
		return err
	}

	guildConfig.AutoRoleIDs = strings.Join(roleIDs, ";")

	return f.SetGuildConfig(guildID, guildConfig)

}

// GetGuildAutoDelete returns the auto delete setting for a guild.
func (f *MiddleWare) GetGuildAutoDelete(guildID string) (bool, error) {

	guildConfig, err := f.GetGuildConfig(guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get guild config")
		return false, err
	}

	return guildConfig.AutoDelete, nil

}

// SetGuildAutoDelete sets the auto delete setting for a guild.
func (f *MiddleWare) SetGuildAutoDelete(guildID string, autoDelete bool) error {

	guildConfig, err := f.GetGuildConfig(guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get guild config")
		return err
	}

	guildConfig.AutoDelete = autoDelete

	return f.SetGuildConfig(guildID, guildConfig)

}

// GuildEntryExists checks if a guild exists.
func (f *MiddleWare) GuildEntryExists(guildID string) (bool, error) {
	guildPath := fmt.Sprintf("%s/guilds/%s.json", f.fio.FIOPath, guildID)

	return f.fio.CheckFile(guildPath), nil
}

// CreateGuildEntry creates a guild entry.
func (f *MiddleWare) CreateGuildEntry(guildID string) error {
	guildPath := fmt.Sprintf("%s/guilds/%s.json", f.fio.FIOPath, guildID)

	// generate the file
	// don't know if we should check before, w/e
	err := f.fio.GenerateFile(guildPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to create guild entry")
		return err
	}

	guildConfig := map[string]interface{}{
		"autoDelete":  false,
		"autoRoleIDs": "",
		"guildID":     guildID,
	}

	return f.fio.Save(guildPath, guildConfig)

}

// DeleteGuildEntry deletes a guild entry.
func (f *MiddleWare) DeleteGuildEntry(guildID string) error {
	guildPath := fmt.Sprintf("%s/guilds/%s.json", f.fio.FIOPath, guildID)

	err := f.fio.DeleteFile(guildPath)
	if err != nil {
		logrus.WithField("guildID", guildID).WithError(err).Error("Failed to delete guild config")
		return err
	}

	return err
}

// GetUserConfig returns the user config for a user.
func (f *MiddleWare) GetUserConfig(userID string) (models.UserConfig, error) {
	userPath := fmt.Sprintf("%s/users/%s.json", f.fio.FIOPath, userID)

	userMap, err := f.fio.Parse(userPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse user config")
		return models.UserConfig{}, err
	}

	return models.UserConfig{
		UserID: userMap["user_id"].(string),
	}, err
}

// SetUserConfig sets the user config for a user.
func (f *MiddleWare) SetUserConfig(userID string, userConfig models.UserConfig) error {
	userPath := fmt.Sprintf("%s/users/%s.json", f.fio.FIOPath, userID)

	userMap := map[string]interface{}{
		"user_id": userConfig.UserID,
	}

	return f.fio.Save(userPath, userMap)
}

// UserEntryExists checks if a user exists.
func (f *MiddleWare) UserEntryExists(userID string) (bool, error) {
	userPath := fmt.Sprintf("%s/users/%s.json", f.fio.FIOPath, userID)

	return f.fio.CheckFile(userPath), nil
}

// CreateUserEntry creates a user entry.
func (f *MiddleWare) CreateUserEntry(userID string) error {
	userPath := fmt.Sprintf("%s/users/%s.json", f.fio.FIOPath, userID)

	err := f.fio.GenerateFile(userPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to create user entry")
		return err
	}

	userMap := map[string]interface{}{
		"user_id": userID,
	}

	return f.fio.Save(userPath, userMap)
}

// DeleteUserEntry deletes a user entry.
func (f *MiddleWare) DeleteUserEntry(userID string) error {
	userPath := fmt.Sprintf("%s/users/%s.json", f.fio.FIOPath, userID)

	err := f.fio.DeleteFile(userPath)
	if err != nil {
		logrus.WithField("userID", userID).WithError(err).Error("Failed to delete user config")
		return err
	}

	return err

}

// GetUserXPEntry returns the user level entry for a user.
func (f *MiddleWare) GetUserXPEntry(userID, guildID string) (models.UserXPEntry, error) {
	xpPath := fmt.Sprintf("%s/xp/%s/%s.json", f.fio.FIOPath, guildID, userID)

	xpMap, err := f.fio.Parse(xpPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse user level entry")
		return models.UserXPEntry{}, err
	}

	return models.UserXPEntry{
		UserID:               xpMap["user_id"].(string),
		Level:                int(xpMap["level"].(float64)),
		CurrentXP:            int(xpMap["current_xp"].(float64)),
		TotalXP:              int(xpMap["total_xp"].(float64)),
		LastMessageTimestamp: int64(xpMap["last_message_timestamp"].(float64)),
		LastSessionID:        xpMap["last_session_id"].(string),
		LastSessionTimestamp: int64(xpMap["last_session_timestamp"].(float64)),
	}, err
}

// SetUserXPEntry sets the user level entry for a user.
func (f *MiddleWare) SetUserXPEntry(userID, guildID string, userXPEntry models.UserXPEntry) error {
	xpPath := fmt.Sprintf("%s/xp/%s/%s.json", f.fio.FIOPath, guildID, userID)

	xpMap := map[string]interface{}{
		"user_id":           userXPEntry.UserID,
		"level":             userXPEntry.Level,
		"current_xp":        userXPEntry.CurrentXP,
		"total_xp":          userXPEntry.TotalXP,
		"last_message":      userXPEntry.LastMessageTimestamp,
		"last_session_id":   userXPEntry.LastSessionID,
		"last_session_time": userXPEntry.LastSessionTimestamp,
	}

	return f.fio.Save(xpPath, xpMap)
}

// GetUserLevel returns the user level for a user.
func (f *MiddleWare) GetUserLevel(userID, guildID string) (int, error) {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse user xp entry")
		return 0, err
	}

	return xpEntry.Level, nil
}

// SetUserLevel sets the user level for a user.
func (f *MiddleWare) SetUserLevel(userID, guildID string, level int) error {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return err
	}

	xpEntry.Level = level

	return f.SetUserXPEntry(userID, guildID, xpEntry)

}

// GetUserCurrentXP returns the user current xp for a user.
func (f *MiddleWare) GetUserCurrentXP(userID, guildID string) (int, error) {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return 0, err
	}

	return xpEntry.CurrentXP, nil
}

// SetUserCurrentXP sets the user current xp for a user.
func (f *MiddleWare) SetUserCurrentXP(userID, guildID string, currentXP int) error {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return err
	}

	xpEntry.CurrentXP = currentXP

	return f.SetUserXPEntry(userID, guildID, xpEntry)
}

// GetUserTotalXP returns the user total xp for a user.
func (f *MiddleWare) GetUserTotalXP(userID, guildID string) (int, error) {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return 0, err
	}

	return xpEntry.TotalXP, nil
}

// SetUserTotalXP sets the user total xp for a user.
func (f *MiddleWare) SetUserTotalXP(userID, guildID string, totalXP int) error {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return err
	}

	xpEntry.TotalXP = totalXP

	return f.SetUserXPEntry(userID, guildID, xpEntry)
}

// GetUserLastMessageTimestamp returns the last message timestamp for a user.
func (f *MiddleWare) GetUserLastMessageTimestamp(userID, guildID string) (int64, error) {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)

	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return 0, err
	}

	return xpEntry.LastMessageTimestamp, nil
}

// SetUserLastMessageTimestamp sets the last message timestamp for a user.
func (f *MiddleWare) SetUserLastMessageTimestamp(userID, guildID string, lastMessageTimestamp int64) error {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return err
	}

	xpEntry.LastMessageTimestamp = lastMessageTimestamp

	return f.SetUserXPEntry(userID, guildID, xpEntry)
}

// GetUserLastSessionID returns the last session id for a user.
func (f *MiddleWare) GetUserLastSessionID(userID, guildID string) (string, error) {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return "", err
	}

	return xpEntry.LastSessionID, nil
}

// SetUserLastSessionID sets the last session id for a user.
func (f *MiddleWare) SetUserLastSessionID(userID, guildID string, lastSessionID string) error {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return err
	}

	xpEntry.LastSessionID = lastSessionID

	return f.SetUserXPEntry(userID, guildID, xpEntry)
}

// GetUserLastSessionTimestamp returns the last session timestamp for a user.
func (f *MiddleWare) GetUserLastSessionTimestamp(userID, guildID string) (int64, error) {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return 0, err
	}

	return xpEntry.LastSessionTimestamp, nil
}

// SetUserLastSessionTimestamp sets the last session timestamp for a user.
func (f *MiddleWare) SetUserLastSessionTimestamp(userID, guildID string, lastSessionTimestamp int64) error {
	xpEntry, err := f.GetUserXPEntry(userID, guildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return err
	}

	xpEntry.LastSessionTimestamp = lastSessionTimestamp

	return f.SetUserXPEntry(userID, guildID, xpEntry)
}

// UserXPEntryExists checks if a user level entry exists.
func (f *MiddleWare) UserXPEntryExists(userID, guildID string) (bool, error) {
	xpPath := fmt.Sprintf("%s/xp/%s/%s.json", f.fio.FIOPath, guildID, userID)

	return f.fio.CheckFile(xpPath), nil
}

// CreateUserXPEntry creates a user level entry.
func (f *MiddleWare) CreateUserXPEntry(userID, guildID string) error {
	xpPath := fmt.Sprintf("%s/xp/%s/%s.json", f.fio.FIOPath, guildID, userID)

	err := f.fio.GenerateFile(xpPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to create user xp entry")
		return err
	}

	xpMap := map[string]interface{}{
		"user_id":                userID,
		"level":                  0,
		"current_xp":             0,
		"total_xp":               0,
		"last_message_timestamp": 0,
		"last_session_id":        "",
		"last_session_timestamp": 0,
	}

	return f.fio.Save(xpPath, xpMap)
}

// DeleteUserXPEntry deletes a user level entry.
func (f *MiddleWare) DeleteUserXPEntry(userID, guildID string) error {
	xpPath := fmt.Sprintf("%s/xp/%s/%s.json", f.fio.FIOPath, guildID, userID)

	return f.fio.DeleteFile(xpPath)
}
