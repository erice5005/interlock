package interlock

import (
	"fmt"
)

type Node interface {
	// Port    int
	// Active  bool
	Link() error
	Run()
	WriteTo(string, MSGType, interface{})
	Read() <-chan DataFrame
	// SetType(nk NodeKind)
	SetConfig(nc NodeConfig)
	Close()
}

type NodeConfig struct {
	ID      string
	Address string
	Port    int
}

func (nc NodeConfig) ToString() string {
	return fmt.Sprintf("%v:%v", nc.Address, nc.Port)
}

type LinkOp int64

const (
	Connect LinkOp = 0
	Serve   LinkOp = 1
)

type NodeKind int64

const (
	Routing      NodeKind = 0
	Fanout       NodeKind = 1
	PubSub       NodeKind = 2
	NodeClient   NodeKind = 3
	PubsubClient NodeKind = 4
)

func NewNode(kx NodeKind) Node {
	switch kx {
	case Routing:
		return NewRouter()
	case NodeClient:
		return NewClient()
	}
	return nil
}
