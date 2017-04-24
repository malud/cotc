package main

import (
	"fmt"
	"testing"
)

func TestSearchSolution(t *testing.T) {
	ships, mines, barrels, balls := NewTestData()
	state := NewState(balls, mines, barrels, ships)
	solution := SearchSolution(state)
	fmt.Printf("%+v", solution)
}
