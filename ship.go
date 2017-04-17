package main

const (
	ShipLength       = 3
	ShipWidth        = 1
	ShipMaxRum       = 100
	ShipCannonRange  = 10
	ShipFireTurnLock = 1
)

type Ship struct {
	Entity
	target   Coord
	rotation int
	speed    int
	rum      int
	owner    int
}

func (s *Ship) Dist(o *Ship) int {
	return s.pos.DistanceTo(o.pos)
}
