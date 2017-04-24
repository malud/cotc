package main

import (
	"fmt"
	"math"
)

const (
	ShipLength         = 3
	ShipWidth          = 1
	ShipMaxRum         = 100
	ShipCannonRange    = 10
	ShipCannonCooldown = 1
	ShipMineCooldown   = 4
	ShipMinSpeed       = 0
	ShipMaxSpeed       = 2
)

type Ship struct {
	Entity
	target             Coord
	action             string
	orientation        int
	speed              int
	rum                int
	initialRum         int
	owner              int
	mineCooldown       int
	cannonCooldown     int
	newOrientation     int
	newPosition        *Coord
	newBowCoordinate   *Coord
	newSternCoordinate *Coord
}

func (s *Ship) Stern() Coord {
	return s.pos.Neighbor((s.orientation + 3) % 6)
}

func (s *Ship) Bow() Coord {
	return s.pos.Neighbor(s.orientation)
}

func (s *Ship) NewStern() Coord {
	return s.pos.Neighbor((s.newOrientation + 3) % 6)
}

func (s *Ship) NewBow() Coord {
	return s.pos.Neighbor(s.newOrientation)
}

func (s *Ship) At(pos Coord) bool {
	stern := s.Stern()
	bow := s.Bow()
	return stern.Equals(pos) || bow.Equals(pos) || s.pos.Equals(pos)
}

func (s *Ship) Damage(health int) {
	s.rum -= health
	if s.rum <= 0 {
		s.rum = 0
	}
}

func (s *Ship) Heal(health int) {
	s.rum += health
	if s.rum > ShipMaxRum {
		s.rum = ShipMaxRum
	}
}

func (s *Ship) NewBowIntersect(other *Ship) bool {
	return s.newBowCoordinate != nil && (s.newBowCoordinate.Equals(*other.newBowCoordinate) || s.newBowCoordinate.Equals(*other.newPosition) || s.newBowCoordinate.Equals(*other.newSternCoordinate))
}

func (s *Ship) NewBowIntersectAll(ships []*Ship) bool {
	for _, other := range ships {
		if s != other && s.NewBowIntersect(other) {
			return true
		}
	}
	return false
}

func (s *Ship) NewPositionsIntersect(other *Ship) bool {
	sternCollision := s.newSternCoordinate != nil && (s.newSternCoordinate.Equals(*other.newBowCoordinate) || s.newSternCoordinate.Equals(*other.newPosition) || s.newSternCoordinate.Equals(*other.newSternCoordinate))
	centerCollision := s.newPosition != nil && (s.newPosition.Equals(*other.newBowCoordinate) || s.newPosition.Equals(*other.newPosition) || s.newPosition.Equals(*other.newSternCoordinate))
	return s.NewBowIntersect(other) || sternCollision || centerCollision
}

func (s *Ship) NewPositionsIntersectAll(ships []*Ship) bool {
	for _, other := range ships {
		if s != other && s.NewPositionsIntersect(other) {
			return true
		}
	}
	return false
}

func (s *Ship) Dist(o *Ship) float64 {
	return float64(s.pos.DistanceTo(o.pos))
}

func (s *Ship) String() string {
	return fmt.Sprintln(s.id, s.rum, s.speed, s.pos)
}
