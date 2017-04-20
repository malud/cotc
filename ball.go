package main

type Ball struct {
	Entity
	from       *Ship
	target     *Ship
	travelTime int
}
