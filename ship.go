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

func (s *Ship) MoveTo(targetPosition Coord) {
	currentPosition := s.pos
	if currentPosition.Equals(targetPosition) {
		s.action = ActionSlower
	}

	var targetAngle, angleStraight, anglePort, angleStarboard, centerAngle, anglePortCenter, angleStarboardCenter float64

	switch s.speed {
	case ShipMaxSpeed:
		s.action = ActionSlower
	case 1:
		// Suppose we've moved first
		currentPosition = currentPosition.Neighbor(s.orientation)
		if !currentPosition.IsInsideMap() {
			s.action = ActionSlower
		}

		// Target reached at next turn
		if currentPosition.Equals(targetPosition) {
			s.action = ""
		}

		// For each neighbor cell, find the closest to target
		targetAngle = currentPosition.Angle(targetPosition)
		angleStraight = math.Min(math.Abs(float64(s.orientation)-targetAngle), 6-math.Abs(float64(s.orientation)-targetAngle))
		anglePort = math.Min(math.Abs((float64(s.orientation)+1)-targetAngle), math.Abs((float64(s.orientation)-5)-targetAngle))
		angleStarboard = math.Min(math.Abs((float64(s.orientation)+5)-targetAngle), math.Abs((float64(s.orientation)-1)-targetAngle))

		centerAngle = currentPosition.Angle(Coord{MapWidth / 2, MapHeight / 2})
		anglePortCenter = math.Min(math.Abs((float64(s.orientation)+1)-centerAngle), math.Abs((float64(s.orientation)-5)-centerAngle))
		angleStarboardCenter = math.Min(math.Abs((float64(s.orientation)+5)-centerAngle), math.Abs((float64(s.orientation)-1)-centerAngle))

		// Next to target with bad angle, slow down then rotate (avoid to turn around the target!)
		if currentPosition.DistanceTo(targetPosition) == 1 && angleStraight > 1.5 {
			s.action = ActionSlower
		}

		distanceMin := -1

		// Test forward
		nextPosition := currentPosition.Neighbor(s.orientation)
		if nextPosition.IsInsideMap() {
			distanceMin = nextPosition.DistanceTo(targetPosition)
			s.action = ""
		}

		// Test port
		nextPosition = currentPosition.Neighbor((s.orientation + 1) % 6)
		if nextPosition.IsInsideMap() {
			distance := nextPosition.DistanceTo(targetPosition)
			if distanceMin == -1 || distance < distanceMin || distance == distanceMin && anglePort < angleStraight-0.5 {
				distanceMin = distance
				s.action = ActionPort
			}
		}

		// Test starboard
		nextPosition = currentPosition.Neighbor((s.orientation + 5) % 6)
		if nextPosition.IsInsideMap() {
			distance := nextPosition.DistanceTo(targetPosition)
			if distanceMin == -1 || distance < distanceMin || (distance == distanceMin && angleStarboard < anglePort-0.5 && s.action == ActionPort) || (distance == distanceMin && angleStarboard < angleStraight-0.5 && s.action == "") || (distance == distanceMin && s.action == ActionPort && angleStarboard == anglePort && angleStarboardCenter < anglePortCenter) || (distance == distanceMin && s.action == ActionPort && angleStarboard == anglePort && angleStarboardCenter == anglePortCenter && (s.orientation == 1 || s.orientation == 4)) {
				distanceMin = distance
				s.action = ActionStarBoard
			}
		}
	case 0:
		// Rotate ship towards target
		targetAngle = currentPosition.Angle(targetPosition)
		angleStraight = math.Min(math.Abs(float64(s.orientation)-targetAngle), 6-math.Abs(float64(s.orientation)-targetAngle))
		anglePort = math.Min(math.Abs((float64(s.orientation)+1)-targetAngle), math.Abs((float64(s.orientation)-5)-targetAngle))
		angleStarboard = math.Min(math.Abs((float64(s.orientation)+5)-targetAngle), math.Abs((float64(s.orientation)-1)-targetAngle))

		centerAngle = currentPosition.Angle(Coord{MapWidth / 2, MapHeight / 2})
		anglePortCenter = math.Min(math.Abs((float64(s.orientation)+1)-centerAngle), math.Abs((float64(s.orientation)-5)-centerAngle))
		angleStarboardCenter = math.Min(math.Abs((float64(s.orientation)+5)-centerAngle), math.Abs((float64(s.orientation)-1)-centerAngle))

		forwardPosition := currentPosition.Neighbor(s.orientation)

		s.action = ""

		if anglePort <= angleStarboard {
			s.action = ActionPort
		}

		if angleStarboard < anglePort || angleStarboard == anglePort && angleStarboardCenter < anglePortCenter || angleStarboard == anglePort && angleStarboardCenter == anglePortCenter && (s.orientation == 1 || s.orientation == 4) {
			s.action = ActionStarBoard
		}

		if forwardPosition.IsInsideMap() && angleStraight <= anglePort && angleStraight <= angleStarboard {
			s.action = ActionFaster
		}
	}
	if s.action == "" {
		s.action = ActionFaster
	}
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

//func (s *Ship) MoveTo(dst Coord) string {
//	angle := s.pos.Angle(dst)
//	n := s.pos.Neighbor(int(angle))
//	return s.RotateTo(dst)
//}

func (s *Ship) RotateTo(dst Coord) string {
	angle := s.pos.Angle(dst)
	if angle < float64(s.orientation) {
		return ActionStarBoard // right
	} else if angle > float64(s.orientation) {
		return ActionPort // left
	}
	return ActionWait
}

func (s *Ship) Dist(o *Ship) float64 {
	return float64(s.pos.DistanceTo(o.pos))
}

func (s *Ship) String() string {
	return fmt.Sprintln(s.id, s.rum, s.speed, s.pos)
}
