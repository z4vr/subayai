package xp

import "math"

// NeededXP returns the amount of XP needed to reach the next level
func (u *UserXP) NeededXP() int {
	return 5*int(math.Pow(float64(u.Level), 2)) + (50 * u.Level) + 100 - u.CurrentXP
}

// AddXP adds XP to the user recursively until there is
// no more XP to add
func (u *UserXP) AddXP(earnedXP int) {
	neededXP := u.NeededXP()
	if u.CurrentXP >= neededXP {
		u.Level++
		earnedXP -= neededXP
		u.CurrentXP = 0
		u.AddXP(earnedXP)
	} else {
		u.CurrentXP += earnedXP
	}
}
