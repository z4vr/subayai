package leveling

type Data struct {
	UserID    string
	GuildID   string
	CurrentXP int
	TotalXP   int
	Level     int
}

type GuildData map[string]UserData

type UserData map[string]*Data
