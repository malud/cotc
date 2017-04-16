package main

const (
	TypeShip       = "SHIP"
	TypeBarrel     = "BARREL"
	TypeCannonBall = "CANNONBALL"
	TypeMine       = "MINE"
)

type Entity struct {
	id  int
	pos Vec
}
