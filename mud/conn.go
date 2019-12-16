package mud

import (
	"bufio"
	"net"
	"strings"
	"time"
)

type Conn struct {
	conn     *net.TCPConn
	SendChan chan []byte
	RecvChan chan []byte
	Quit     chan bool
}

func (c *Conn) GetID() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) Close() {
	c.conn.Close()
	c.Quit <- true
}

func (c *Conn) Send(message []byte) error {
	err := c.conn.SetWriteDeadline(time.Now().Add(time.Second * 3))
	if err != nil {
		c.Close()
		return err
	}
	_, err = c.conn.Write(message)
	if err != nil {
		c.Close()
	}
	return err
}

func (c *Conn) loop() {
	defer c.conn.Close()
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.Quit <- true
			//log.Printf("conn err %v", err)
			break
		}
		msg := strings.TrimRight(message, "\r\n")
		//log.Printf("Read msg: %s", msg)
		c.RecvChan <- []byte(msg)
	}
}

func NewConn(conn *net.TCPConn) *Conn {
	c := &Conn{
		conn:     conn,
		SendChan: make(chan []byte),
		RecvChan: make(chan []byte),
		Quit:     make(chan bool, 2),
	}
	go c.loop()
	return c
}
