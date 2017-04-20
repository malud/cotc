package main

type Player struct {
	id         int
	ships      []*Ship
	shipsAlive []*Ship
}

func (p *Player) SetDead() {
	for _, ship := range p.ships {
		ship.rum = 0
	}
}

func (p *Player) GetScore() int {
	var rum int
	for _, ship := range p.ships {
		rum += ship.rum
	}
	return rum
}
