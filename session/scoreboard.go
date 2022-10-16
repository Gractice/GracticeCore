package session

type Scoreboard func(*Player)

func (s Scoreboard) Send(session *Player) {
	s(session)
}
