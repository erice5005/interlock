package interlock

import (
	"net"
	"testing"
	"time"
)

func Test_Conn(t *testing.T) {
	server, client := net.Pipe()

	n1 := NewConnection(server)
	n2 := NewConnection(client)

	n1.Run()
	n2.Run()
	t.Logf("N2 ID: %v\n", n2.ID)
	n2.Send(CONNECT, nil)
	time.Sleep(10 * time.Millisecond)
}

func Test_ConnErr(t *testing.T) {
	server, client := net.Pipe()

	n1 := NewConnection(server)
	n2 := NewConnection(client)

	n1.Run()
	n2.Run()

	n2.Send(CONNECT, nil)

	go func() {
		time.Sleep(3 * time.Second)
		n1.Close()
	}()

	time.Sleep(10 * time.Second)
}
