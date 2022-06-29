package interlock

import (
	"log"
	"net"
	"time"
)

type Router struct {
	Connections map[string]*Connection
	Config      NodeConfig
	InputStream chan DataFrame
	LX          net.Listener
	end         chan bool
}

func NewRouter() *Router {
	return &Router{
		Connections: make(map[string]*Connection),
		InputStream: make(chan DataFrame),
		end:         make(chan bool),
	}
}

func (r *Router) SetConfig(nc NodeConfig) {
	r.Config = nc
}

func (r *Router) Link() error {
	l, err := net.Listen("tcp", r.Config.ToString())
	if err != nil {
		return err
	}
	r.LX = l
	return nil
}

func (r *Router) Run() {
	defer r.LX.Close()
	defer func() {
		log.Printf("End Run")
	}()
	go func() {
		for range time.Tick(500 * time.Millisecond) {
			for cid, cnn := range r.Connections {
				if cnn.closed {
					cnn.Close()
					delete(r.Connections, cid)
				}
			}
		}
	}()
	ncc, ecc := r.acceptLoop()
	for {
		select {
		case <-r.end:
			log.Printf("got end\n")
			for _, cnn := range r.Connections {
				cnn.Close()
			}
			return
		case ex := <-ecc:
			log.Printf("accept err: %v\n", ex)
			r.Close()
		case c := <-ncc:
			newConn := NewConnection(c)
			log.Printf("New Connection\n")
			go newConn.Run()
			go func() {
				for df := range newConn.Input {
					r.InputStream <- df
				}
			}()
			r.Connections[newConn.ID] = newConn
		}
	}

}

func (r *Router) acceptLoop() (<-chan net.Conn, <-chan error) {
	ch := make(chan (net.Conn))
	eCh := make(chan error)

	go func() {
		for {
			c, err := r.LX.Accept()
			if err != nil {
				eCh <- err
				return
			}
			ch <- c
		}
	}()

	return ch, eCh
}

func (r *Router) Read() <-chan DataFrame {
	return r.InputStream
}

func (r *Router) WriteTo(tg string, kx MSGType, data interface{}) {
	if tg != "" {
		if cx, ok := r.Connections[tg]; ok {
			cx.Send(kx, data)
		}
	} else {
		for _, cx := range r.Connections {
			cx.Send(kx, data)
		}
	}
}

func (r *Router) Close() {
	close(r.end)
	r.LX.Close()
}
