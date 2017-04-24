package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	Enemy int = iota
	Me
)

func debugln(pattern string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", pattern), args...)
}

func main() {
	lastRoundFired := make([]bool, 0)
	for {
		var myShipCount int
		fmt.Scan(&myShipCount)

		var entityCount int
		fmt.Scan(&entityCount)

		balls := make([]*Ball, 0)
		barrels := make([]*Barrel, 0)
		ships := make([]*Ship, 0)
		mines := make([]*Mine, 0)
		//fleet := NewFleet() // just my ships for now

		for i := 0; i < entityCount; i++ {
			var entityId int
			var entityType string
			var x, y, arg1, arg2, arg3, arg4 int
			fmt.Scan(&entityId, &entityType, &x, &y, &arg1, &arg2, &arg3, &arg4)
			//debugln("Entity: %d %s %d,%d - %d %d %d %d\n", entityId, entityType, x, y, arg1, arg2, arg3, arg4)
			switch entityType {
			case TypeShip:
				ship := &Ship{
					Entity:      Entity{id: entityId, pos: Coord{x, y}},
					orientation: arg1,
					speed:       arg2,
					rum:         arg3,
					owner:       arg4,
				}
				ships = append(ships, ship)
				if ship.owner == Me {
					//fleet.Join(ship)
					lastRoundFired = append(lastRoundFired, false)
				}
			case TypeBarrel:
				barrel := &Barrel{
					Entity: Entity{id: entityId, pos: Coord{x, y}},
					rum:    arg1,
				}
				barrels = append(barrels, barrel)
			case TypeMine:
				mine := &Mine{
					Entity: Entity{id: entityId, pos: Coord{x, y}},
				}
				mines = append(mines, mine)
			case TypeCannonBall:
				ball := &Ball{
					Entity:     Entity{id: entityId, pos: Coord{x, y}},
					travelTime: arg1,
				}
				balls = append(balls, ball)
			}

		}

		state := NewState(balls, mines, barrels, ships)
		solution := SearchSolution(state)

		fmt.Printf("%s\n", strings.Join(solution.actions, "\n"))
	}
}
