package mud

import (
	"fmt"
	"log"
	"sync"
)

type World struct {
	Messages  chan *Message
	PlayerMap sync.Map
	Count int
}

func (w *World) Welcome(conn *Conn) {
	log.Printf("%v connect", conn.GetRemoteAddr())
	conn.SendChan <- []byte(WelcomeText)
}

func (w *World) Login(conn *Conn) {
	w.Welcome(conn)
	conn.SendChan <- []byte("Your name: ")

	select {
	case name := <-conn.RecvChan:
		player := NewPlayer(w.Messages, conn, conn.GetRemoteAddr(), string(name))
		w.PlayerMap.Store(player.GetID(), player)

		w.Count += 1
		log.Printf("%s %s join the game!", conn.GetRemoteAddr(), string(name))
	}
}

func (w *World) Logout(id string) {
	if player, exist := w.PlayerMap.Load(id); exist {
		log.Printf("%s %s leave the game!", id, player.(*Player).Name)
		w.PlayerMap.Delete(id)
		w.Count -= 1
	}
}

func (w *World) dispatcher(msg *Message) {
	switch msg.Type {
	case LeaveMsg:
		id := msg.Value.(string)
		w.Logout(id)
	case SystemInfoMsg:
		id := msg.Value.(string)
		w.systemInfo(id)
	}
}

func (w *World) broadcast() {

}

func (w *World) systemInfo(playerID string) {
	if player, ok := w.PlayerMap.Load(playerID); ok {
		player.(*Player).Conn.SendChan <- []byte(fmt.Sprintf("World : total player %d\n", w.Count))
	}
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
		//for {
		//	time.Sleep(time.Second * 3)
		//	w.PlayerMap.Range(func(key, value interface{}) bool {
		//		player := value.(*Player)
		//		player.Conn.SendChan <- []byte(fmt.Sprintf("World : total player %d\n", w.Count))
		//		return true
		//	})
		//}
	}()
}

func NewWorld() *World {
	return &World{
		Messages:  make(chan *Message, 10),
		PlayerMap: sync.Map{},
	}
}
