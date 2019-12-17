package mud

import (
	"bufio"
	"net"
	"strings"
)

type Conn struct {
	conn     *net.TCPConn
	SendChan chan []byte
	RecvChan chan []byte
	Quit     chan bool
	close	chan bool
}

func (c *Conn) GetRemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) Close() {
	c.conn.Close()
	c.Quit <- true
	c.close <- true
}

func (c *Conn) loop() {
	defer c.conn.Close()
	// read
	go func() {
		reader := bufio.NewReader(c.conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				c.Close()
				//log.Printf("conn err %v", err)
				break
			}
			msg := strings.TrimRight(message, "\r\n")
			//log.Printf("Read msg: %s", msg)
			c.RecvChan <- []byte(msg)
		}
	}()
	// send
	go func() {
		for {
			select {
			case buff := <-c.SendChan:
				_, err := c.conn.Write(buff)
				if err != nil {
					c.Close()
					break
				}
			}
		}
	}()
	<-c.close
}

func NewConn(conn *net.TCPConn) *Conn {
	c := &Conn{
		conn:     conn,
		SendChan: make(chan []byte, 10),
		RecvChan: make(chan []byte, 10),
		Quit:     make(chan bool, 3),
		close: make(chan bool, 3),
	}
	go c.loop()
	return c
}
