package main

import (
	"fmt"
	"math"
)

const (
	DirectionEven = [][]int{{1, 0}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}}
	DirectionOdd  = [][]int{{1, 0}, {1, -1}, {0, -1}, {-1, 0}, {0, 1}, {1, 1}}
	Directions    = [][]int{{1, -1, 0}, {+1, 0, -1}, {0, +1, -1}, {-1, +1, 0}, {-1, 0, +1}, {0, -1, +1}}
)

type Coord struct {
	x int
	y int
}

type CubeCoord struct {
	x int
	y int
	z int
}

func (c Coord) Angle(targetPosition Coord) float64 {
	dy := (targetPosition.y - c.y) * math.Sqrt(3) / 2
	dx := targetPosition.x - c.x + ((c.y-targetPosition.y)&1)*0.5
	angle := -math.Atan2(dy, dx) * 3 / math.Pi
	if angle < 0 {
		angle += 6
	} else if angle >= 6 {
		angle -= 6
	}
	return angle
}

func (c Coord) ToCubeCoord() CubeCoord {
	xp := c.x - (c.y-(c.y&1))/2
	zp := c.y
	yp := -(xp + zp)
	return CubeCoord{xp, yp, zp}
}

func (c Coord) Neighbor(orientation int) Coord {
	var newY, newX int
	if c.y%2 == 1 {
		newY = c.y + DirectionOdd[orientation][1]
		newX = c.x + DirectionOdd[orientation][0]
	} else {
		newY = c.y + DirectionEven[orientation][1]
		newX = c.x + DirectionEven[orientation][0]
	}

	return Coord{newX, newY}
}

func (c Coord) IsInsideMap() bool {
	return c.x >= 0 && c.x < MapWidth && c.y >= 0 && c.y < MapHeight
}

func (a Coord) DistanceTo(dst Coord) int {
	return a.ToCubeCoord().DistanceTo(dst.ToCubeCoord())
}

func (a Coord) Equals(b Coord) bool {
	return a.y == b.y && a.x == b.x
}

func (c Coord) String() string {
	return fmt.Sprint(c.x, c.y)
}

func (c CubeCoord) ToOffsetCoord() Coord {
	newX := c.x + (c.z-(c.z&1))/2
	newY := c.z
	return Coord{newX, newY}
}

func (c CubeCoord) Neighbor(orientation int) CubeCoord {
	nx := c.x + Directions[orientation][0]
	ny := c.y + Directions[orientation][1]
	nz := c.z + Directions[orientation][2]
	return CubeCoord{nx, ny, nz}
}

func (c CubeCoord) DistanceTo(dst CubeCoord) int {
	return (abs(c.x-dst.x) + abs(c.y-dst.y) + abs(c.z-dst.z)) / 2
}

func (c CubeCoord) String() string {
	return fmt.Sprint(c.x, c.y, c.z)
}

// math.Abs implementation for integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}
