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
	rotation int
	speed    int
	rum      int
	owner    int
}
