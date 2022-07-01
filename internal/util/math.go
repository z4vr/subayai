package util

import "math"

// NeededXP returns the XP needed for the current level and current XP
func NeededXP(level, currentXP int) int {
	return 5*int(math.Pow(float64(level), 2)) + (50 * level) + 100 - currentXP
}

// CurrentLevel returns the current level and xp of a user
// after the earned XP is added to the current XP
func CurrentLevel(earnedXP int, currentLevel, currentXP int) (int, int) {
	if currentXP >= NeededXP(currentLevel, currentXP) {
		return CurrentLevel(earnedXP-NeededXP(currentLevel, currentXP), currentLevel+1, 0)
	} else {
		return currentXP + earnedXP, currentLevel
	}
}
