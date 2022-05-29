package leveling

type LevelData struct {
	UserID    string
	GuildID   string
	CurrentXP int
	TotalXP   int
	Level     int
}

type GuildMap map[string]MemberMap

type MemberMap map[string]*LevelData
