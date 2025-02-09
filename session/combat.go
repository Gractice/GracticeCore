package session

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type CombatInfo struct {
	// Force is the horizon knockback
	Force float64
	// Height is the vertical knockback
	Height float64
}

type combatHandler struct {
	p *Player
}

func (h *combatHandler) HandleAttackEntity(_ *event.Context[*player.Player], _ world.Entity, force, height *float64, _ *bool) {
	if info := h.p.CombatInfo(); info != nil {
		*force, *height = info.Force, info.Height
	}
}
