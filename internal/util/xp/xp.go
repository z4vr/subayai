package xp

import (
	"github.com/z4vr/subayai/internal/models"
	"math"
)

// NeededXP returns the amount of XP needed to reach the next level
func NeededXP(u *models.UserXP) int {
	return 5*int(math.Pow(float64(u.Level), 2)) + (50 * u.Level) + 100 - u.CurrentXP
}

// AddXP adds XP to the user recursively until there is
// no more XP to add
func AddXP(u *models.UserXP, earnedXP int, levelup bool) bool {
	neededXP := NeededXP(u)
	if u.CurrentXP >= neededXP {
		u.Level++
		levelup = true
		earnedXP -= neededXP
		u.CurrentXP = 0
		return AddXP(u, earnedXP, levelup)
	} else {
		u.CurrentXP += earnedXP
		return levelup
	}
}
