package main

import (
	"github.com/threkk/sprx/proxy"
	"log"
	"net"
)

func main() {
	server := proxy.NewProxy()
	listener, _ := net.Listen("tcp", "localhost:0")
	log.Println(listener.Addr())
	log.Fatal(server.Serve(listener))
	// Format: sprx -p 5200 -link "some.corporate.web:22 > :8022, localhost:8080 > :8080" alberto@myip.local
	// 1. Start SSH connection. No connection, exit.
	// 2. Start proxy server. No proxy, exit.
	// 3. Configure PAC.
	// 4. Start all the tunnels, starting with the proxy and following with the
	// PAC and the to the links.
	// 5. Display the information:
	//   -> PAC url (localhost:port)
	//   -> Ports forwarded.
	//   -> Close with Ctrl + C, Ctrl + D
}
