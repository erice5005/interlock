package interlock

import (
	"log"
	"net"
	"time"
)

type Client struct {
	Config     NodeConfig
	Conn       Connection
	end        chan bool
	connects   int
	retryLimit int
}

func NewClient() *Client {
	return &Client{
		end:        make(chan bool),
		connects:   0,
		retryLimit: 3,
	}
}

func (c *Client) SetConfig(nc NodeConfig) {
	c.Config = nc
}

func (c *Client) Link() error {
	cx, err := net.Dial("tcp", c.Config.ToString())
	if err != nil {
		return err
	}

	c.Conn = *NewConnection(cx)
	c.Conn.Send(CONNECT, nil)
	c.connects++
	return nil
}

func (c *Client) Run() {
	defer c.Conn.Close()

	for c.connects < c.retryLimit {
		c.Conn.Run()
		time.Sleep(time.Duration(c.retryLimit) * 500 * time.Millisecond)
		log.Printf("Retry\n")
		err := c.Link()
		if err != nil {
			log.Printf("Link err: %v\n", err)
			return
		}

	}
	// <-c.end
}

func (c *Client) Read() <-chan DataFrame {
	return c.Conn.Input
}

func (c *Client) WriteTo(tg string, kx MSGType, data interface{}) {
	c.Conn.Send(kx, data)
}

func (c *Client) Close() {
	c.end <- true
}
