package session

import (
	"github.com/google/uuid"
	"sync"
)

type Query[T any] func() T

type Data[T any] struct {
	mu      *sync.RWMutex
	entries map[uuid.UUID]T
}

func (d *Data[T]) Query(p *Player) Query[T] {
	return func() T {
		d.mu.RLock()
		defer d.mu.RUnlock()
		return d.entries[p.Id()]
	}
}

func (d *Data[T]) Register(p *Player, data T) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries[p.Id()] = data
}

func (d *Data[T]) UnRegister(p *Player) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, p.Id())
}

func NewData[T any]() *Data[T] {
	return &Data[T]{
		mu:      &sync.RWMutex{},
		entries: make(map[uuid.UUID]T),
	}
}
