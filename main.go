package main

import (
	"SimpleMud/mud"
	"SimpleMud/telnet"
	"log"
	"net"
)

func main() {
	w := mud.NewWorld()
	go w.Start()

	s := telnet.NewServer( 9888, func(conn *net.TCPConn) {
		w.Login(mud.NewConn(conn))
	})
	log.Println("Server start")
	log.Fatal(s.Run())
}
