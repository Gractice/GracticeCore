package arena

import "github.com/blackjack200/gracticecore/util"

const (
	startBitOffset = iota
	closeBitOffset
	canJoinBitOffset
)

type BitState util.BitMap

func (b *BitState) IsStarted() bool {
	return (*util.BitMap)(b).Get(startBitOffset)
}

func (b *BitState) IsClosed() bool {
	return (*util.BitMap)(b).Get(closeBitOffset)
}

func (b *BitState) IsJoinable() bool {
	return (*util.BitMap)(b).Get(canJoinBitOffset)
}

func (b *BitState) Start() {
	(*util.BitMap)(b).On(startBitOffset)
	(*util.BitMap)(b).Off(closeBitOffset)
}

func (b *BitState) CanJoin() {
	(*util.BitMap)(b).On(canJoinBitOffset)
}

func (b *BitState) CannotJoin() {
	(*util.BitMap)(b).Off(canJoinBitOffset)
}

func (b *BitState) Close() {
	(*util.BitMap)(b).On(closeBitOffset)
	(*util.BitMap)(b).Off(startBitOffset)
}

func NewBitState() *BitState {
	return (*BitState)(util.NewBitMap())
}
