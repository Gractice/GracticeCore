package session

import (
	"github.com/df-mc/dragonfly/server/player"
	"golang.org/x/exp/slices"
	"sync"
)

var mu = &sync.Mutex{}
var sessions []*Player

func Lookup(p *player.Player) *Player {
	mu.Lock()
	defer mu.Unlock()
	for _, s := range sessions {
		if s.player == p {
			return s
		}
	}
	panic("this should not happen")
}

func register(p *Player) {
	mu.Lock()
	sessions = append(sessions, p)
	mu.Unlock()
}

func unregister(p *Player) {
	idx := slices.Index(sessions, p)
	if idx != -1 {
		sessions = slices.Delete(sessions, idx, idx)
	}
	mu.Unlock()
}
