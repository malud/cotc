package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	Enemy int = iota
	Me
)

func debugln(pattern string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", pattern), args...)
}

func main() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	lastRoundFired := make([]bool, 0)
	for {
		var myShipCount int
		fmt.Scan(&myShipCount)

		var entityCount int
		fmt.Scan(&entityCount)

		barrels := make([]*Barrel, 0)
		ships := make([]*Ship, 0)
		mines := make([]*Mine, 0)
		actions := make([]string, 0)
		fleet := NewFleet() // just my ships for now

		for i := 0; i < entityCount; i++ {
			var entityId int
			var entityType string
			var x, y, arg1, arg2, arg3, arg4 int
			fmt.Scan(&entityId, &entityType, &x, &y, &arg1, &arg2, &arg3, &arg4)
			//debugln("Entity: %d %s %d,%d - %d %d %d %d\n", entityId, entityType, x, y, arg1, arg2, arg3, arg4)
			switch entityType {
			case TypeShip:
				ship := &Ship{
					Entity:   Entity{id: entityId, pos: Coord{x, y}},
					rotation: arg1,
					speed:    arg2,
					rum:      arg3,
					owner:    arg4,
				}
				ships = append(ships, ship)
				if ship.owner == Me {
					fleet.Join(ship)
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
			}
		}

		myTargets := make([]Coord, 0) // prevent double actions
		for i := 0; i < myShipCount; i++ {
			myShip := fleet.ships[i]

			if !lastRoundFired[i] {
				for _, ship := range ships {
					if ship.owner == Enemy && ship.speed == 0 {
						actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionFire, ship.pos.x, ship.pos.y))
						lastRoundFired[i] = true
						break
					}
				}
				if lastRoundFired[i] {
					continue
				}
			}

			if fleet.ships[i].rum < 100 && len(barrels) > 1 { // shoot at last one
				var nearest = barrels[0]
				for _, barrel := range barrels {
					if barrel.pos.DistanceTo(myShip.pos) < nearest.pos.DistanceTo(myShip.pos) {
						isTarget := false // and barrel isn't already a target
						for _, t := range myTargets {
							if t.x == barrel.pos.x && t.y == barrel.pos.y {
								isTarget = true
								break
							}
						}
						if !isTarget {
							nearest = barrel
						}
					}
				}
				actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionMove, nearest.pos.x, nearest.pos.y))
				lastRoundFired[i] = false
			} else if len(barrels) == 1 && !lastRoundFired[i] {
				actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionFire, barrels[0].pos.x, barrels[0].pos.y))
				lastRoundFired[i] = true
			} else {
				// just move
				x := rnd.Intn(MapWidth)
				y := rnd.Intn(MapHeight)
				actions = append(actions, fmt.Sprintf("%s %d %d\n", ActionMove, x, y))
				lastRoundFired[i] = false
			}
		}
		for _, action := range actions {
			fmt.Print(action)
		}
	}
}
