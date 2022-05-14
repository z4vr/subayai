package models

import "math"

// GuildConfig is a struct for the guild configuration which is stored as a json file
// in the filesystem.
type GuildConfig struct {
	GuildID      string `json:"guild_id"`
	AutoRoleIDs  string `json:"auto_role_ids"`
	AutoDelete   bool   `json:"auto_delete"`
	AFKChannelID string `json:"afk_channel_id"`
	// TODO: add more fields as needed
}

// UserConfig is a struct for the user configuration which is stored as a json file
// in the filesystem.
type UserConfig struct {
	UserID string `json:"user_id"`
	// TODO: add more fields as needed
}

// UserXPEntry is a struct for the user level entry which is stored as a json file
// in the filesystem.
type UserXPEntry struct {
	UserID               string `json:"user_id"`
	Level                int    `json:"level"`
	CurrentXP            int    `json:"current_xp"`
	TotalXP              int    `json:"total_xp"`
	LastMessageTimestamp int64  `json:"last_message_timestamp"`
	LastSessionID        string `json:"last_session_id"`
	LastSessionTimestamp int64  `json:"last_session_timestamp"`
}

func (xp *UserXPEntry) XPNeeded() int {
	return 5*int(math.Pow(float64(xp.Level), 2)) + (50 * xp.Level) + 100 - xp.CurrentXP
}

func (xp *UserXPEntry) LevelUp(earnedXP int) {
	neededXP := xp.XPNeeded()
	if earnedXP >= neededXP {
		xp.Level++
		earnedXP -= neededXP
		xp.CurrentXP = 0
		xp.LevelUp(earnedXP)
	} else {
		xp.CurrentXP += earnedXP
	}
}
