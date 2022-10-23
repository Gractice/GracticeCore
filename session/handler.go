package session

import (
	"github.com/Blackjack200/GracticeEssential/mhandler"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/gractice/gracticecore/arena"
	"github.com/gractice/gracticecore/util"
	"time"
)

type LifecycleHandler interface {
	Close()
}

type Player struct {
	closed  *atomic.Bool
	player  *player.Player
	handler *mhandler.MultipleHandler
	ticker  *time.Ticker

	scoreboard     *atomic.Value[Scoreboard]
	currentHandler *atomic.Value[any]
	arena          *atomic.Value[arena.Arena]
}

func NewSession(player *player.Player) *Player {
	p := &Player{
		closed:         atomic.NewBool(false),
		player:         player,
		handler:        mhandler.New(),
		ticker:         time.NewTicker(time.Second / 20),
		scoreboard:     atomic.NewValue[Scoreboard](nil),
		currentHandler: atomic.NewValue[any](nil),
		arena:          atomic.NewValue[arena.Arena](nil),
	}
	player.Handle(p.handler)
	// OnQuit
	p.handler.Register(p)
	go p.tick()
	register(p)
	return p
}

func (p *Player) Id() uuid.UUID {
	return p.RealPlayer().UUID()
}

func (p *Player) tick() {
	for !p.closed.Load() {
		<-p.ticker.C
		v := p.scoreboard.Load()
		if v != nil {
			v.Send(p)
		}
	}
}

func (p *Player) RealPlayer() *player.Player {
	return p.player
}

func (p *Player) SwitchWorld(w *world.World) {
	w.AddEntity(p.RealPlayer())
	p.Reset(true)
	p.RealPlayer().Teleport(w.Spawn().Vec3Middle())
}

func (p *Player) Reset(teleport bool) {
	pp := p.player

	if teleport {
		pp.Teleport(pp.World().PlayerSpawn(pp.UUID()).Vec3Middle())
	}
	pp.Heal(pp.MaxHealth()-pp.Health(), healing.SourceInstantHealthEffect{})
	//hardcoded
	pp.SetFood(20)
	p.Clear()
}

func (p *Player) Clear() {
	pp := p.player

	pp.Inventory().Clear()
	pp.EnderChestInventory().Clear()
	for _, e := range pp.Effects() {
		pp.RemoveEffect(e.Type())
	}
	pp.Extinguish()
	pp.AbortBreaking()
	pp.ResetFallDistance()
	pp.SetExperienceLevel(0)
	pp.SetExperienceProgress(0)
}

func (p *Player) CurrentHandler() any {
	return p.currentHandler.Load()
}

func (p *Player) Handle(h any) (unregisterFunc func()) {
	old := p.currentHandler.Swap(h)
	p.handler.Unregister(old)
	util.Cast[LifecycleHandler](old).Is(func(h LifecycleHandler) {
		h.Close()
	})

	p.handler.Register(h)
	return func() {
		p.handler.Unregister(h)
		p.currentHandler.Store(nil)
		util.Cast[LifecycleHandler](h).Is(func(h LifecycleHandler) {
			h.Close()
		})
	}
}

func (p *Player) HandleGlobal(h any) (unregisterFunc func()) {
	p.handler.Register(h)
	return func() {
		p.handler.Unregister(h)
		util.Cast[LifecycleHandler](h).Is(func(h LifecycleHandler) {
			h.Close()
		})
	}
}

func (p *Player) SendMessage(msg string, args ...any) {
	p.RealPlayer().Messagef(msg, args...)
}

func (p *Player) SetScoreboard(sb Scoreboard) {
	if sb == nil {
		p.RealPlayer().RemoveScoreboard()
	}
	p.scoreboard.Swap(sb)
	sb.Send(p)
}

func (p *Player) Scoreboard() Scoreboard {
	return p.scoreboard.Load()
}

func (p *Player) Kill() {
	p.player.Hurt(p.player.MaxHealth(), damage.SourceVoid{})
}

func (p *Player) Arena() arena.Arena {
	return p.arena.Load()
}

func (p *Player) LeaveArena() {
	oldArena := p.arena.Swap(nil)
	if oldArena != nil {
		_ = oldArena.Remove(p)
	}
}

func (p *Player) JoinArena(a arena.Arena) error {
	p.LeaveArena()
	return a.Add(p)
}

func (p *Player) OnJoin(a arena.Arena) {
	p.arena.Store(a)
}

func (p *Player) OnQuit(arena.Arena) {
	p.arena.Store(nil)
}

func (p *Player) close() {
	p.closed.Store(true)
	p.LeaveArena()
	p.scoreboard.Store(nil)
	p.handler.Clear()
	p.ticker.Stop()
	mu.Lock()
	unregister(p)
}

func (p *Player) HandleQuit() {
	p.close()
}
