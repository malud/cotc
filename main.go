package main

import (
	"fmt"
	"os"
)

const (
	Enemy int = iota
	Me
)

func main() {
	//grid := &Grid{
	//	width: Width,
	//	height: Height,
	//}

	var lastRoundFired bool
	for {
		// myShipCount: the number of remaining ships
		var myShipCount int
		fmt.Scan(&myShipCount)

		// entityCount: the number of entities (e.g. ships, mines or cannonballs)
		var entityCount int
		fmt.Scan(&entityCount)

		barrels := make([]*Barrel, 0)
		ships := make([]*Ship, 0)
		mines := make([]*Mine, 0)
		actions := make([]string, 0)

		for i := 0; i < entityCount; i++ {
			var entityId int
			var entityType string
			var x, y, arg1, arg2, arg3, arg4 int
			fmt.Scan(&entityId, &entityType, &x, &y, &arg1, &arg2, &arg3, &arg4)
			//fmt.Fprintf(os.Stderr, "Entity: %d %s %d,%d - %d %d %d %d\n", entityId, entityType, x, y, arg1, arg2, arg3, arg4)
			switch entityType {
			case TypeShip:
				ship := &Ship{
					Entity:   Entity{id: entityId, position: Vec{x, y}},
					rotation: arg1,
					speed:    arg2,
					rum:      arg3,
					owner:    arg4,
				}
				ships = append(ships, ship)
			case TypeBarrel:
				barrel := &Barrel{
					Entity: Entity{id: entityId, position: Vec{x, y}},
					rum:    arg1,
				}
				barrels = append(barrels, barrel)
			case TypeMine:
				mine := &Mine{
					Entity: Entity{id: entityId, position: Vec{x, y}},
				}
				mines = append(mines, mine)
			}
		}
		for i := 0; i < myShipCount; i++ {
			var myShip *Ship
			for _, ship := range ships {
				if ship.owner == Me {
					myShip = ship
				}
			}

			var lastDist int
			for _, ship := range ships {
				if ship.owner != Me {
					lastDist = dist(myShip.position, ship.position)
					fmt.Fprintf(os.Stderr, "Dist: %d", dist(myShip.position, ship.position))
				}
			}

			if lastDist > 2 && myShip.rum < 100 && len(barrels) > 0 {
				var nearest = barrels[0]
				for _, barrel := range barrels {
					if dist(barrel.position, myShip.position) < dist(nearest.position, myShip.position) {
						nearest = barrel
					}
				}
				fmt.Fprintf(os.Stderr, "Barrel: %#v\n", nearest)
				actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionMove, nearest.position.x, nearest.position.y))
			} else {
				for _, ship := range ships {
					if ship.owner == Enemy {
						if !lastRoundFired && dist(myShip.position, ship.position) <= ShipCannonRange {
							fmt.Fprintf(os.Stderr, "Fire at %+v\n", ship)
							actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionFire, ship.position.x, ship.position.y))
							lastRoundFired = true
						} else {
							actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionMove, ship.position.x, ship.position.y))
							lastRoundFired = false
						}
					}
				}
			}

			// fmt.Fprintln(os.Stderr, "Debug messages...")
			//fmt.Printf("MOVE 11 10\n") // Any valid action, such as "WAIT" or "MOVE x y"
		}
		for _, action := range actions {
			fmt.Print(action)
		}
	}
}

func dist(a Vec, b Vec) int {
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
