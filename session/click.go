package session

import (
	"github.com/Blackjack200/GracticeEssential/mhandler"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"sync"
	"time"
)

type clickHandler struct {
	mu           *sync.Mutex
	clickCounter uint32
	lastClick    time.Time
	OnClick      func()
}

var _ mhandler.PunchAirHandler = (*clickHandler)(nil)

func (h *clickHandler) HandlePunchAir(ctx *event.Context[*player.Player]) {
	h.click()
}

var _ mhandler.AttackEntityHandler = (*clickHandler)(nil)

func (h *clickHandler) HandleAttackEntity(_ *event.Context[*player.Player], _ world.Entity, _, _ *float64, _ *bool) {
	h.click()
}

func (h *clickHandler) click() {
	mu.Lock()
	if time.Since(h.lastClick) > time.Second {
		h.clickCounter = 0
	}
	h.clickCounter++
	h.lastClick = time.Now()
	mu.Unlock()
	if h.OnClick != nil {
		h.OnClick()
	}
}

func (h *clickHandler) CPS() uint32 {
	mu.Lock()
	defer mu.Unlock()
	if time.Since(h.lastClick) > time.Second {
		h.clickCounter = 0
	}
	return h.clickCounter
}
