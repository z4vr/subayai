package level

import (
	"math"
)

// NeededXP returns the amount of XP needed to reach the next level
func (x *Data) NeededXP() int {
	return 5*int(math.Pow(float64(x.Level), 2)) + (50 * x.Level) + 100 - x.CurrentXP
}

// AddXP adds XP to the user recursively until there is
// no more XP to add
func (x *Data) AddXP(earnedXP int, levelup bool) bool {
	neededXP := x.NeededXP()
	if x.CurrentXP >= neededXP {
		x.Level++
		levelup = true
		earnedXP -= neededXP
		x.CurrentXP = 0
		return x.AddXP(earnedXP, levelup)
	} else {
		x.CurrentXP += earnedXP
		return levelup
	}
}
