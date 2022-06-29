package interlock

import (
	"bufio"
	"io"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

type ConnData struct {
	Conn     net.Conn
	FrameIn  chan DataFrame
	FrameOut chan DataFrame
	close    chan bool
	active   bool
	closed   chan bool
}

func newByteReaderStream(r io.Reader, size int) (<-chan []byte, <-chan error) {
	if size <= 0 {
		size = 2048
	}

	ch := make(chan ([]byte))
	errChan := make(chan error)
	buf := bufio.NewReader(r)
	go func() {
		// for {

		// inBuf := make([]byte, size)
		// s := 0
		// inner:
		for {
			message, err := buf.ReadString('\n')
			// n, err := r.Rea
			// n, err := r.Read(inBuf[s:])
			// if n > 0 {
			// 	ch <- inBuf[s : s+n]
			// 	s += n
			// }
			if err != nil {
				log.Printf("Read Err: %v\n", err)
				errChan <- err
				return
			}
			ch <- []byte(message)
			// if s >= len(inBuf) {
			// 	break inner
			// }
		}
		// }
	}()
	return ch, errChan
}

func NewConnData(c net.Conn) ConnData {
	return ConnData{
		Conn:     c,
		FrameIn:  make(chan DataFrame, 1),
		FrameOut: make(chan DataFrame, 1),
		close:    make(chan bool),
		closed:   make(chan bool),
	}
}

func (cd *ConnData) Deactivate() {
	cd.active = false
}

func (cd *ConnData) Loop() {
	cd.active = true
	defer cd.Deactivate()
	byteStream, errStream := newByteReaderStream(cd.Conn, 4096)
	for {
		select {
		case <-cd.close:
			log.Printf("close\n")
			cd.closed <- true
			return
		case toSend := <-cd.FrameOut:
			cd.Conn.Write(toSend.ToBytes())
		case readErr := <-errStream:
			log.Printf("Read err: %v\n", readErr)
			cd.Close()
		case inp := <-byteStream:
			if inp == nil {
				continue
			}
			recv := *ParseDataFrame([]byte(inp))
			cd.FrameIn <- recv
		}
	}
}

func (cd *ConnData) Close() {
	if !cd.active {
		return
	}
	cd.Deactivate()
	close(cd.FrameIn)
	close(cd.FrameOut)
	cd.Conn.Close()
	close(cd.close)
	// cd.Conn = nil
	log.Printf("Closed ConnData\n")

}

type Connection struct {
	ID           string
	LinkID       string
	Conn         ConnData
	Input        chan DataFrame
	Output       chan DataFrame
	PingInterval time.Duration
	LastPing     time.Time
	log          ConnectionLog
	closed       bool
}

func (c *Connection) read() {

	for recv := range c.Conn.FrameIn {
		c.log.Add(recv)
		switch recv.Kind {
		case PING:
			c.LastPing = time.Now()
			// log.Printf("[LOG] %v got ping from %v at %v\n", c.ID, c.LinkID, c.LastPing)
		case CONNECT:
			c.LastPing = time.Now()
			c.LinkID = recv.SrcID
			log.Printf("[LOG] %v got connect from %v\n", c.ID, c.LinkID)
		case DATA:
			c.Input <- recv
		case DISCONNECT:
			log.Printf("[LOG] %v got disconnect\n", c.ID)
			c.Close()
		}
	}
}

func (c *Connection) write() {
	for toSend := range c.Output {
		// c.Conn.Write(toSend.ToBytes())
		if c.closed {
			return
		}
		c.Conn.FrameOut <- toSend
	}
}

func (c *Connection) Run() {
	// go
	defer c.Close()
	go c.Conn.Loop()
	go c.read()
	go c.write()
	<-c.Conn.closed
	c.closed = true
	log.Printf("End run loop\n")
	// go c.Ping()
}

func (c *Connection) Close() {
	if c.closed {
		return
	}
	c.closed = true
	// if _, ok := <-c.Input; ok {
	// 	close(c.Input)
	// }
	// if _, ok := <-c.Output; ok {

	// }
	if c.Conn.active {
		// c.Send(DISCONNECT, nil)
		c.Conn.Close()
	}
	close(c.Input)
	close(c.Output)
}

func (c *Connection) Ping() {
	if c.PingInterval == 0 {
		return
	}
	pingTick := time.NewTicker(c.PingInterval)
	for range pingTick.C {
		if c.closed {
			return
		}
		c.Output <- DataFrame{Kind: PING}
	}
}

func NewConnection(c net.Conn) *Connection {
	return &Connection{
		ID:           uuid.New().String(),
		Conn:         NewConnData(c),
		Input:        make(chan DataFrame, 1),
		Output:       make(chan DataFrame, 1),
		PingInterval: 500 * time.Millisecond,
		log:          NewLog(),
	}
}

func (c *Connection) Send(kx MSGType, dx interface{}) {
	c.Output <- DataFrame{
		Kind:    kx,
		Content: dx,
		Created: time.Now(),
		SrcID:   c.ID,
	}
}
