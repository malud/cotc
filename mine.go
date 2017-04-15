package main

// a ship cannot place another for the next 4 turns
const MineTurnLock = 4

type Mine struct {
	Entity
	r int
}
