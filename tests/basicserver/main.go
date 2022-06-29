package main

import (
	"fmt"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/erice5005/interlock"
)

// type RouterServer struct {
// 	Count int
// 	RX    *interlock.RouterNode
// }

// func (rs *RouterServer) Do(df *interlock.DataFrame) {
// 	rs.Count++
// 	log.Printf("Server Got: %v\n", df)
// 	rs.RX.SendTo(df.SrcID, &interlock.DataFrame{
// 		Kind:    interlock.DATA,
// 		Content: df.Content,
// 	})
// }

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	sx := interlock.NewNode(interlock.Routing)
	sx.SetConfig(interlock.NodeConfig{
		Address: "localhost",
		Port:    10000,
	})
	sx.Link()
	go func() {
		for msg := range sx.Read() {
			log.Printf("SX Got: %v, %v\n", msg, msg.Content)
		}
	}()
	// go func() {
	// 	time.Sleep(30 * time.Second)
	// 	log.Printf("Close SX\n")
	// 	sx.Close()
	// }()
	sx.Run()
}
