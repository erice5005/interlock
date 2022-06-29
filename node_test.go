package interlock

// type EchoHandler struct {
// 	Prefix string
// }

// func (eh EchoHandler) Do(df *DataFrame) {
// 	log.Printf("%v Got: %v", eh.Prefix, df)
// }

// func NewEchoHandler(prefix string) EchoHandler {
// 	return EchoHandler{
// 		Prefix: prefix,
// 	}
// }

// func Test_ComLineRuns(t *testing.T) {
// 	server, client := net.Pipe()

// 	comServer := NewCommLine(server)
// 	comClient := NewCommLine(client)

// 	comServer.Handler = NewEchoHandler("server")
// 	comClient.Handler = NewEchoHandler("client")

// 	go comServer.Run()
// 	go comClient.Run()

// 	comServer.Output <- DataFrame{
// 		Kind: PING,
// 	}

// 	comClient.Output <- DataFrame{
// 		Kind: PONG,
// 	}

// 	time.Sleep(5 * time.Second)
// }
