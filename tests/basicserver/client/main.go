package main

import (
	"log"
	"time"

	_ "net/http/pprof"

	"github.com/erice5005/interlock"
)

func main() {
	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	// nx.Handler =
	nx := interlock.NewNode(interlock.NodeClient)
	nx.SetConfig(interlock.NodeConfig{
		Address: "localhost",
		Port:    10000,
	})
	nx.Link()
	go func() {
		// time.Sleep(10 * time.Second)
		// nx.WriteTo("", interlock.DISCONNECT, nil)
		// for {
		// 	reader := bufio.NewReader(os.Stdin)
		// 	fmt.Print(">> ")
		// 	text, _ := reader.ReadString('\n')
		// 	nx.WriteTo("", interlock.DATA, text)
		// }
		for range time.NewTicker(1 * time.Second).C {
			nx.WriteTo("", interlock.DATA, "Hello")
		}
	}()
	log.Printf("Start Client\n")
	nx.Run()
	log.Printf("end client\n")
}
