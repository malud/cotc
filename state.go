package main

import (
	"math"
)

const (
	LowDamage            = 25
	HighDamage           = 50
	MineDamage           = 25
	RewardRumBarrelValue = 30
)

//private long seed;
//private int shipsPerPlayer;
//private int mineCount;
//private int barrelCount;
//private Random random;
type State struct {
	cannonballs []*Ball
	mines       []*Mine
	barrels     []*Barrel
	ships       []*Ship
	damage      []Damage
	explosions  []Coord // cannonBallExplosions
	players     []*Player
}

// Create a new State to simulate the next turn and create copies for all references
func NewState(balls []*Ball, mines []*Mine, barrels []*Barrel, ships []*Ship) *State {
	newBalls := make([]*Ball, 0)
	for _, ball := range balls {
		b := *ball
		newBalls = append(newBalls, &b)
	}
	newMines := make([]*Mine, 0)
	for _, mine := range mines {
		m := *mine
		newMines = append(newMines, &m)
	}
	newBarrels := make([]*Barrel, 0)
	for _, barrel := range barrels {
		b := *barrel
		newBarrels = append(newBarrels, &b)
	}
	newShips := make([]*Ship, 0)
	for _, ship := range ships {
		s := *ship
		newShips = append(newShips, &s)
	}
	state := &State{
		cannonballs: newBalls,
		mines:       newMines,
		barrels:     newBarrels,
		ships:       newShips,
		damage:      make([]Damage, 0),
		explosions:  make([]Coord, 0),
		players:     make([]*Player, 0),
	}

	for i := 0; i < 2; i++ {
		playerShips := make([]*Ship, 0)
		for _, s := range newShips {
			if s.owner == i {
				playerShips = append(playerShips, s)
			}
		}
		p := &Player{
			id:         i,
			ships:      playerShips,
			shipsAlive: playerShips,
		}
		state.players = append(state.players, p)
	}
	return state
}

func (s *State) Clone() *State {
	return NewState(s.cannonballs, s.mines, s.barrels, s.ships)
}

func (s *State) Update() {
	s.moveCannonballs()
	s.decrementRum()
	s.updateInitialRum()

	s.applyActions()
	s.moveShips()
	s.rotateShips()

	s.explodeShips()
	s.explodeMines()
	s.explodeBarrels()

	// For each sunk ship, create a new rum barrel with the amount of rum the ship had at the begin of the turn (up to 30).
	for _, ship := range s.ships {
		if ship.rum <= 0 {
			reward := min(RewardRumBarrelValue, ship.initialRum)
			if reward > 0 {
				s.barrels = append(s.barrels, &Barrel{Entity: Entity{pos: ship.pos}, rum: reward})
			}
		}
	}

	// cannonball ui effect - obsolete here
	//for _, position := range s.explosions {
	//	s.damage = append(s.damage, Damage{position, 0, false})
	//}

	for _, ship := range s.ships {
		if ship.rum <= 0 {
			for _, player := range s.players {
				alive := player.shipsAlive[:0]
				for _, s := range player.shipsAlive {
					if s != ship {
						alive = append(alive, s)
					}
				}
				player.shipsAlive = alive
			}
		}
	}

	//if (s.gameIsOver()) {
	//	throw new GameOverException("endReached");
	//}
}

func (s *State) moveCannonballs() {
	for _, cannonball := range s.cannonballs {
		if cannonball.travelTime > 0 {
			cannonball.travelTime--
		} else if cannonball.travelTime == 0 {
			s.explosions = append(s.explosions, cannonball.pos)
		}
	}
}

func (s *State) decrementRum() {
	for _, s := range s.ships {
		s.Damage(1)
	}
}

func (s *State) updateInitialRum() {
	for _, s := range s.ships {
		s.initialRum = s.rum
	}
}

func (s *State) applyActions() {
	for _, player := range s.players {
		for _, ship := range player.shipsAlive {
			if ship.mineCooldown > 0 {
				ship.mineCooldown--
			}
			if ship.cannonCooldown > 0 {
				ship.cannonCooldown--
			}

			ship.newOrientation = ship.orientation

			if ship.action != "" {
				switch ship.action {
				case ActionFaster:
					if ship.speed < ShipMaxSpeed {
						ship.speed++
					}
				case ActionSlower:
					if ship.speed > 0 {
						ship.speed--
					}
				case ActionPort:
					ship.newOrientation = (ship.orientation + 1) % 6
				case ActionStarBoard:
					ship.newOrientation = (ship.orientation + 5) % 6
				case ActionMine:
					if ship.mineCooldown == 0 {
						target := ship.Stern().Neighbor((ship.orientation + 3) % 6)
						if target.IsInsideMap() {
							cellIsFreeOfBarrels, cellIsFreeOfMines, cellIsFreeOfShips := true, true, true
							for _, b := range s.barrels {
								if b.pos.Equals(target) {
									cellIsFreeOfBarrels = false
									break
								}
							}
							for _, m := range s.mines {
								if m.pos.Equals(target) {
									cellIsFreeOfMines = false
									break
								}
							}
							for _, s := range s.ships {
								if s.pos.Equals(target) {
									cellIsFreeOfShips = false
									break
								}
							}

							if cellIsFreeOfBarrels && cellIsFreeOfShips && cellIsFreeOfMines {
								ship.mineCooldown = ShipMineCooldown
								mine := &Mine{Entity: Entity{pos: target}}
								s.mines = append(s.mines, mine)
							}
						}

					}
				case ActionFire:
					distance := ship.Bow().DistanceTo(ship.target)
					if ship.target.IsInsideMap() && distance <= ShipCannonRange && ship.cannonCooldown == 0 {
						travelTime := int(1 + math.Trunc(float64(ship.Bow().DistanceTo(ship.target)/3.0)))
						s.cannonballs = append(s.cannonballs, &Ball{Entity: Entity{pos: ship.target}, target: ship, from: ship, travelTime: travelTime})
						ship.cannonCooldown = ShipCannonCooldown
					}
					break
				default:
					break
				}
			}
		}
	}
}

func (s *State) moveShips() {
	// ---
	// Go forward
	// ---
	for i := 1; i <= ShipMaxSpeed; i++ {
		for _, player := range s.players {
			for _, ship := range player.shipsAlive {
				a := ship.pos
				ship.newPosition = &a
				b := ship.Bow()
				ship.newBowCoordinate = &b
				c := ship.Stern()
				ship.newSternCoordinate = &c

				if i > ship.speed {
					continue
				}

				newCoordinate := ship.pos.Neighbor(ship.orientation)

				if newCoordinate.IsInsideMap() {
					// Set new coordinate.
					ship.newPosition = &newCoordinate
					a := newCoordinate.Neighbor(ship.orientation)
					ship.newBowCoordinate = &a
					b := newCoordinate.Neighbor((ship.orientation + 3) % 6)
					ship.newSternCoordinate = &b
				} else {
					// Stop ship!
					ship.speed = ShipMinSpeed
				}
			}
		}

		// Check ship and obstacles collisions
		collisions := make([]*Ship, 0)
		collisionDetected := true
		for collisionDetected {
			collisionDetected = false

			for _, ship := range s.ships {
				if ship.NewBowIntersectAll(s.ships) {
					collisions = append(collisions, ship)
				}
			}

			for _, ship := range collisions {
				// Revert last move
				a := ship.pos
				ship.newPosition = &a
				b := ship.Bow()
				ship.newBowCoordinate = &b
				c := ship.Stern()
				ship.newSternCoordinate = &c

				// Stop ships
				ship.speed = ShipMinSpeed

				collisionDetected = true
			}
			collisions = []*Ship{}
		}

		// Move ships to their new location
		for _, ship := range s.ships {
			ship.pos = *ship.newPosition
		}

		// Check collisions
		for _, ship := range s.ships {
			s.checkCollisions(ship)
		}
	}
}

func (s *State) rotateShips() {
	// Rotate
	for _, player := range s.players {
		for _, ship := range player.shipsAlive {
			a := ship.pos
			ship.newPosition = &a
			b := ship.NewBow()
			ship.newBowCoordinate = &b
			c := ship.NewStern()
			ship.newSternCoordinate = &c
		}
	}

	// Check collisions
	collisionDetected := true
	collisions := make([]*Ship, 0)
	for collisionDetected {
		collisionDetected = false

		for _, ship := range s.ships {
			if ship.NewPositionsIntersectAll(s.ships) {
				collisions = append(collisions, ship)
			}
		}

		for _, ship := range collisions {
			ship.newOrientation = ship.orientation
			a := ship.NewBow()
			ship.newBowCoordinate = &a
			b := ship.NewStern()
			ship.newSternCoordinate = &b
			ship.speed = ShipMinSpeed
			collisionDetected = true
		}

		collisions = []*Ship{}
	}

	// Apply rotation
	for _, ship := range s.ships {
		ship.orientation = ship.newOrientation
	}

	// Check collisions
	for _, ship := range s.ships {
		s.checkCollisions(ship)
	}
}

func (s *State) checkCollisions(ship *Ship) {
	bow := ship.Bow()
	stern := ship.Stern()
	center := ship.pos

	// Collision with the barrels
	for _, barrel := range s.barrels {
		if barrel.pos.Equals(bow) || barrel.pos.Equals(stern) || barrel.pos.Equals(center) {
			ship.Heal(barrel.rum)
		}
	}

	// Collision with the mines
	for _, mine := range s.mines {
		mineDamage := mine.Explode(s.ships, false)

		if len(mineDamage) > 0 {
			s.damage = append(s.damage, mineDamage...)
		}
	}
}

func (s *State) explodeShips() {
	for _, position := range s.explosions {
		for _, ship := range s.ships {
			if position.Equals(ship.Bow()) || position.Equals(ship.Stern()) {
				s.damage = append(s.damage, Damage{position, LowDamage, true})
				ship.Damage(LowDamage)
				break
			} else if position.Equals(ship.pos) {
				s.damage = append(s.damage, Damage{position, HighDamage, true})
				ship.Damage(HighDamage)
				break
			}
		}
	}
}

func (s *State) explodeMines() {
	for _, position := range s.explosions {
		for _, mine := range s.mines {
			if mine.pos.Equals(position) {
				s.damage = append(s.damage, mine.Explode(s.ships, true)...)
				break
			}
		}
	}
}

func (s *State) explodeBarrels() {
	for _, position := range s.explosions {
		for _, barrel := range s.barrels {
			if barrel.pos.Equals(position) {
				s.damage = append(s.damage, Damage{position, 0, true})
				break
			}
		}
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
