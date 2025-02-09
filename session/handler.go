package session

import (
	"github.com/Blackjack200/GracticeEssential/mhandler"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/gractice/gracticecore/arena"
	"sync"
	"time"
)

type Player struct {
	closed  *atomic.Bool
	player  *player.Player
	handler *mhandler.MultipleHandler
	ticker  *time.Ticker
	click   *clickHandler

	scoreboard         *atomic.Value[Scoreboard]
	currentHandlerFunc *atomic.Value[func()]
	combat             *atomic.Value[*CombatInfo]
	arena              *atomic.Value[arena.Arena]
}

func NewSession(player *player.Player) *Player {
	p := &Player{
		closed:  atomic.NewBool(false),
		player:  player,
		handler: mhandler.New(),
		ticker:  time.NewTicker(time.Second / 20),
		click: &clickHandler{
			mu:           &sync.Mutex{},
			clickCounter: 0,
			lastClick:    time.Now(),
			OnClick:      nil,
		},

		scoreboard:         atomic.NewValue[Scoreboard](nil),
		combat:             atomic.NewValue[*CombatInfo](nil),
		currentHandlerFunc: atomic.NewValue[func()](nil),
		arena:              atomic.NewValue[arena.Arena](nil),
	}

	player.Handle(p.handler)
	p.handler.Register(&combatHandler{p})
	p.handler.Register(p.click)
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
	p.RealPlayer().H().ExecWorld(func(tx *world.Tx, e world.Entity) {
		tx.AddEntity(p.RealPlayer().H())
	})
	p.Reset(true)
	p.RealPlayer().Teleport(w.Spawn().Vec3Middle())
}

func (p *Player) Reset(teleport bool) {
	pp := p.player

	if teleport {
		pp.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
			pp.Teleport(tx.World().PlayerSpawn(pp.UUID()).Vec3Middle())
		})
	}
	pp.Heal(pp.MaxHealth()-pp.Health(), effect.InstantHealingSource{})
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
	return p.currentHandlerFunc.Load()
}

func (p *Player) Handle(h any) (unregisterFunc func()) {
	hh := p.handler.Register(h)
	old := p.currentHandlerFunc.Swap(hh)
	if old != nil {
		old()
	}
	return func() {
		hh()
		p.currentHandlerFunc.Store(nil)
	}
}

func (p *Player) HandleGlobal(h any) (unregisterFunc func()) {
	hh := p.handler.Register(h)
	return func() {
		hh()
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
	p.player.Hurt(p.player.MaxHealth(), entity.VoidDamageSource{})
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

func (p *Player) CombatInfo() *CombatInfo {
	return p.combat.Load()
}

func (p *Player) SetCombatInfo(i *CombatInfo) {
	p.combat.Store(i)
}

func (p *Player) OnClick(f func()) {
	p.click.OnClick = f
}

func (p *Player) CPS() uint32 {
	return p.click.CPS()
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
