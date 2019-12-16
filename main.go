package main

import (
	"SimpleMud/mud"
	"log"
)

func main() {
	w := mud.NewWorld()
	go w.Start()

	s := mud.NewServer(w, 9888)
	log.Fatal(s.Run())
}
