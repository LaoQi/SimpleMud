package mud

type Player struct {
	WorldChan chan *Message
	Conn      *Conn
	Name      string
	ID        string
}

func (p *Player) GetID() string {
	return p.ID
}

func (p *Player) loop() {
	for {
		select {
		case <-p.Conn.Quit:
			p.WorldChan <- NewMessageWithValue(LeaveMsg, p.ID)
			return
		case msg := <-p.Conn.RecvChan:
			switch string(msg) {
			case "quit":
				p.Conn.Close()
			case "system":
				p.WorldChan <- NewMessageWithValue(SystemInfoMsg, p.ID)
			}
		}
	}
}

func NewPlayer(worldChan chan *Message, conn *Conn, id string, name string) *Player {
	p := &Player{
		WorldChan: worldChan,
		Conn:      conn,
		Name:      name,
		ID:        id,
	}
	go p.loop()
	return p
}
