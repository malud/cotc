package main

import "sort"

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

func (f *Fleet) SortByRum() {
	sort.Sort(f)
}

func (f *Fleet) Len() int {
	return len(f.ships)
}

func (f *Fleet) Less(i, j int) bool {
	return f.ships[i].rum < f.ships[j].rum
}

func (f *Fleet) Swap(i, j int) {
	f.ships[i], f.ships[j] = f.ships[j], f.ships[i]
}
