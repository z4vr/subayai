// Package math includes some functions for math, like calculation experience needed for leveling
// and so on.
package math

import "math"

// NeededXP returns the XP needed for the current level to level up
func NeededXP(level int) int {
	return 5*int(math.Pow(float64(level), 2)) + (50 * level) + 100
}

// CurrentLevel returns the current level of the user based on the total XP
// and current level.
func CurrentLevel(currentXP, level int) int {
	neededXP := NeededXP(level)
	if currentXP >= neededXP {
		return CurrentLevel(currentXP-neededXP, level+1)
	} else {
		return level
	}
}
