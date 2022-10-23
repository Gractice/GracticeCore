package arena

type Member interface {
	OnJoin(Arena)
	OnQuit(Arena)
}

type HandleableMember interface {
	Member
	Handle(any) (unReg func())
}

type Arena interface {
	Open() error
	Close() error
	Players() []Member
	Add(Member) error
	Remove(Member) error
}
