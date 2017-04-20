package main

import "fmt"

type Damage struct {
	position Coord
	health   int
	hit      bool
}

func (d *Damage) String() string {
	h := 0
	if d.hit {
		h = 1
	}
	return fmt.Sprint(d.position, d.health, h)
}
