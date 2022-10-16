package arena

type Member interface {
	OnJoin(Arena)
	OnQuit(Arena)
}

type HandleableMember interface {
	Member
	Handle(any) (unReg func())
}

type State interface {
	IsStarted() bool
	IsClosed() bool
	IsJoinable() bool
	Start()
	CanJoin()
	CannotJoin()
	Close()
}

type Descriptor interface {
	Name() string
}

type Arena interface {
	Open() error
	Close() error
	Descriptor() Descriptor
	State() State
	Players() []Member
	Add(Member) error
	Remove(Member) error
}
