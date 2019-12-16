package mud

import (
	"log"
	"sync"
	"time"
)

type World struct {
	Messages  chan *Message
	PlayerMap sync.Map
}

func (w *World) Welcome(conn *Conn) {
	log.Printf("%v connect", conn.GetID())
	_ = conn.Send([]byte(WelcomeText))
}

func (w *World) Login(conn *Conn) {
	w.Welcome(conn)
	_ = conn.Send([]byte("Your name: "))

	select {
	case name := <-conn.RecvChan:
		player := NewPlayer(w.Messages, conn, string(name))
		w.PlayerMap.Store(conn.GetID(), player)

		log.Printf("%s %s join the game!", conn.GetID(), string(name))
	}
}

func (w *World) Logout(id string) {
	if player, exist := w.PlayerMap.Load(id); exist {
		log.Printf("%s %s leave the game!", id, player.(*Player).Name)
		w.PlayerMap.Delete(id)
	}
}

func (w *World) dispatcher(msg *Message) {
	switch msg.Type {
	case LeaveMsg:
		id := msg.Value.(string)
		w.Logout(id)
	}
}

func (w *World) broadcast() {

}

func (w *World) Start() {
	go func() {
		for {
			select {
			case msg := <-w.Messages:
				log.Printf("World Read message %v", msg)
				w.dispatcher(msg)
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 3)
			w.PlayerMap.Range(func(key, value interface{}) bool {
				player := value.(*Player)
				err := player.Conn.Send([]byte("World\n"))
				if err != nil {
					log.Printf("World Broadcast error: %s", err)
					return false
				}
				return true
			})
		}
	}()
}

func NewWorld() *World {
	return &World{
		Messages:  make(chan *Message, 10),
		PlayerMap: sync.Map{},
	}
}
