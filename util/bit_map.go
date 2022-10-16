package util

import "fmt"

func validate(i uint32) {
	if i > 64 {
		panic("bit map just have 64 bit")
	}
}

type BitMap struct {
	data uint64
}

func (b *BitMap) On(i uint32) {
	validate(i)
	b.data |= 1 << i
}

func (b *BitMap) Off(i uint32) {
	validate(i)
	b.data &= ^(1 << i)
}

func (b *BitMap) Get(i uint32) bool {
	validate(i)
	return (b.data & (1 << i)) == 1
}

func (b *BitMap) Clear() {
	b.data &= 0
}

func (b *BitMap) String() string {
	return fmt.Sprintf("%b", b.data)
}

func (b *BitMap) Foreach(fn func(uint32, bool) bool) {
	for i := uint32(0); i < 64; i++ {
		fn(i, (b.data&(1<<i)) == 1)
	}
}

func NewBitMap() *BitMap {
	return &BitMap{}
}
