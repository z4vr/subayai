package xp

type UserXP struct {
	UserID    string
	GuildID   string
	CurrentXP int
	TotalXP   int
	Level     int
}

func New(userID, guildID string, currentXP, totalXP, level int) *UserXP {
	return &UserXP{
		UserID:    userID,
		GuildID:   guildID,
		CurrentXP: currentXP,
		TotalXP:   totalXP,
		Level:     level,
	}
}
