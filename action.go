package main

const (
	ActionFaster = "FASTER"
	ActionFire   = "FIRE"
	ActionMine   = "MINE"
	//ActionMove      = "MOVE"
	ActionSlower    = "SLOWER"
	ActionStarBoard = "STARBOARD"
	ActionPort      = "PORT"
	ActionWait      = "WAIT"
)

var Actions = []string{ActionWait, ActionFaster, ActionSlower, ActionStarBoard, ActionPort, ActionFire, ActionMine}
