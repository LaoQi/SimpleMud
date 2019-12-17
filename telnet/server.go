package telnet

import (
	"log"
	"net"
)

type Server struct {
	Address  *net.TCPAddr
	Listener *net.TCPListener
	AcceptHandler func(conn *net.TCPConn)
}

func (s *Server) loop() error {
	defer s.Listener.Close()
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.AcceptHandler(conn)
	}
}

func (s *Server) Run() error {
	var err error
	s.Listener, err = net.ListenTCP("tcp", s.Address)
	if err != nil {
		return err
	}
	return s.loop()
}

func NewServer(port int, acceptHandler func(conn *net.TCPConn)) *Server {
	return &Server{
		Address: &net.TCPAddr{
			IP:   nil,
			Port: port,
			Zone: "",
		},
		AcceptHandler:acceptHandler,
	}
}
