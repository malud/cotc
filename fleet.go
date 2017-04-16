package main

type Fleet struct {
	ships []*Ship
}

func NewFleet() *Fleet {
	return &Fleet{
		ships: make([]*Ship, 0),
	}
}

func (f *Fleet) Join(s *Ship) {
	f.ships = append(f.ships, s)
}
