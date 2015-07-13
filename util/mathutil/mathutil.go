package mathutil

import (
//"math"
)

/*
	returns the smaller interger
*/
func MinInt(a, b int) int {
	min := b
	if a < b {
		min = a
	}
	return min
}

/*
	returns the absolut value of a
*/
func AbsInt(a int) int {
	if a < 0 {
		return -1 * a
	} else {
		return a
	}
}

/*
	returns the larger interger
*/
func MaxInt(a, b int) int {
	max := b
	if a > b {
		max = a
	}
	return max
}
