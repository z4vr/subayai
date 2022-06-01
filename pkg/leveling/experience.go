package leveling

import (
	"math"
)

// NeededXP returns the amount of XP needed to reach the next leveling
func (d *LevelData) NeededXP() int {
	return 5*int(math.Pow(float64(d.Level), 2)) + (50 * d.Level) + 100 - d.CurrentXP
}

// AddXP adds XP to the user recursively until there is
// no more XP to add
func (d *LevelData) LevelUp(earnedXP int, levelup bool) bool {
	neededXP := d.NeededXP()
	if d.CurrentXP >= neededXP {
		d.Level++
		levelup = true
		earnedXP -= neededXP
		d.CurrentXP = 0
		return d.LevelUp(earnedXP, levelup)
	} else {
		d.CurrentXP += earnedXP
		return levelup
	}
}
