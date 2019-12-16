package mud

type Player struct {
	WorldChan chan *Message
	Conn      *Conn
	Name      string
}

func (p *Player) loop() {
	for {
		select {
		case <-p.Conn.Quit:
			p.WorldChan <- NewMessageWithValue(LeaveMsg, p.Conn.GetID())
			return
		case msg := <-p.Conn.RecvChan:
			switch string(msg) {
			case "quit":
				p.Conn.Close()
			case "system":

			}
		}
	}
}

func NewPlayer(worldChan chan *Message, conn *Conn, name string) *Player {
	p := &Player{
		WorldChan: worldChan,
		Conn:      conn,
		Name:      name,
	}
	go p.loop()
	return p
}
