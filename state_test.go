package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestState_Update(t *testing.T) {
	ships := make([]*Ship, 0)
	ship := &Ship{
		Entity: Entity{
			pos: Coord{MapWidth / 2, MapHeight/2 - 1},
		},
		rum:         ShipMaxRum,
		orientation: 1,
	}
	ships = append(ships, ship)
	state := NewState(make([]*Ball, 0), make([]*Mine, 0), make([]*Barrel, 0), ships)
	state.Update()

	updatedShip := state.ships[0]
	if updatedShip.rum != ShipMaxRum-1 {
		t.Errorf("Expected to lose 1 rum per turn. Got: '%d'", updatedShip.rum)
	}

	state = state.Fork()
	state.ships[0].action = ActionPort
	state.Update()
	if ship.orientation == state.ships[0].orientation || state.ships[0].orientation != 1+1 {
		t.Errorf("Expected to turn left. Got wrong orientation: '%d'", state.ships[0].orientation)
	}

	// forking last state - turning back to orientation 1
	state = state.Fork()
	state.ships[0].action = ActionStarBoard
	state.Update()
	if ship.orientation != state.ships[0].orientation {
		t.Errorf("Expected to turn right. Got wrong orientation: '%d'", state.ships[0].orientation)
	}
}

func NewTestData() ([]*Ship, []*Mine, []*Barrel, []*Ball) {
	src := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(src)
	ships := make([]*Ship, 0)
	mines := make([]*Mine, 0)
	barrels := make([]*Barrel, 0)
	balls := make([]*Ball, 0)

	randCoord := func() Coord {
		return Coord{rand.Intn(MapWidth), rand.Intn(MapHeight)}
	}

	for i := 0; i < 6; i++ {
		ship := &Ship{
			Entity: Entity{id: i, pos: randCoord()},
			owner:  i % 2,
			rum:    ShipMaxRum,
			action: Actions[rand.Intn(len(Actions)-1)],
		}
		ships = append(ships, ship)
	}

	for i := 0; i < 20; i++ {
		mines = append(mines, &Mine{Entity{pos: randCoord()}})
	}

	for i := 0; i < 14; i++ {
		barrels = append(barrels, &Barrel{Entity{pos: randCoord()}, rand.Intn(BarrelMaxRum)})
	}

	for i := 0; i < rand.Intn(12); i++ {
		balls = append(balls, &Ball{Entity: Entity{pos: randCoord()}, travelTime: rand.Intn(ShipCannonRange)})
	}

	return ships, mines, barrels, balls
}

func BenchmarkState_Update(b *testing.B) {
	b.ReportAllocs()
	ships, mines, barrels, balls := NewTestData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		state := NewState(balls, mines, barrels, ships)
		state.Update()
	}
}
