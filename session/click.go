package session

import (
	"github.com/df-mc/dragonfly/server/event"
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

func (h *clickHandler) HandlePunchAir(*event.Context) {
	h.click()
}

func (h *clickHandler) HandleAttackEntity(*event.Context, world.Entity, *float64, *float64, *bool) {
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
