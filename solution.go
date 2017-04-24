package main

import (
	"math/rand"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func SearchSolution(state *State) *Node {
	result := make(chan *Node)
	quit := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(33+3-len(state.ships)) * time.Millisecond)
		quit <- struct{}{}
	}()

	depth := 9 - len(state.ships)
	//solutions := 0

	go func() {
		for {
			select {
			case <-quit:
				close(result)
				return
			default:
				s := state.Clone()
				score := 0.0
				actions := []string{}
				for i := 0; i < depth; i++ {
					es, ea := prepareClone(s)
					s.Update()
					score += s.Score() + es
					if i == 0 {
						actions = append(actions, ea...)
					}
				}
				n := &Node{score: score, actions: actions}
				//solutions++
				result <- n
			}
		}
	}()

	var best *Node
	for n := range result {
		if best == nil {
			best = n
			continue
		}

		if n.score > best.score {
			best = n
		}
	}
	//debugln("Solutions: %d", solutions)
	return best
}

func prepareClone(state *State) (score float64, actions []string) {
	return // :)
}
