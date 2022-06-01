package leveling

type LevelData struct {
	UserID    string
	GuildID   string
	CurrentXP int
	TotalXP   int
	Level     int
}

type DatabaseFunctions interface {
	SetUserCurrentXP(userID, guildID string, currentXP int) error
	GetUserCurrentXP(userID, guildID string) (int, error)
	SetUserTotalXP(userID, guildID string, totalXP int) error
	GetUserTotalXP(userID, guildID string) (int, error)
	SetUserLevel(userID, guildID string, level int) error
	GetUserLevel(userID, guildID string) (int, error)
}
