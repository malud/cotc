package main

type Vec struct {
	x int
	y int
}

func (a Vec) Dist(b Vec) int {
	return max(abs(a.x-b.x), abs(a.y-b.y))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
