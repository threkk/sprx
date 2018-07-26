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
}
