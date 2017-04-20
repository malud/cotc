package main

const NearMineDamage = 10

type Mine struct {
	Entity
}

func (m *Mine) Explode(ships []*Ship, force bool) []Damage {
	damage := make([]Damage, 0)
	var victim *Ship

	for _, ship := range ships {
		if m.pos.Equals(ship.Bow()) || m.pos.Equals(ship.Stern()) || m.pos.Equals(ship.pos) {
			damage = append(damage, Damage{m.pos, MineDamage, true})
			ship.Damage(MineDamage)
			victim = ship
		}
	}

	if force || victim != nil {
		if victim == nil {
			damage = append(damage, Damage{m.pos, MineDamage, true})
		}

		for _, ship := range ships {
			if ship != victim {
				var impactPosition *Coord
				if ship.Stern().DistanceTo(m.pos) <= 1 {
					p := ship.Stern()
					impactPosition = &p
				}
				if ship.Bow().DistanceTo(m.pos) <= 1 {
					p := ship.Bow()
					impactPosition = &p
				}
				if ship.pos.DistanceTo(m.pos) <= 1 {
					impactPosition = &ship.pos
				}

				if impactPosition != nil {
					ship.Damage(NearMineDamage)
					damage = append(damage, Damage{*impactPosition, NearMineDamage, true})
				}
			}
		}
	}

	return damage
}
