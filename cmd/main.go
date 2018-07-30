package main

import (
	"flag"
	"fmt"
	fw "github.com/threkk/sprx/forwarding"
	// "github.com/threkk/sprx/proxy"
	// "log"
	/// "net"
	"os"
)

var port int
var tunnelFlag fw.Tunnels

func usage(out *os.File, description bool) {
	if description {
		// TODO: Add description. Docs too.
		fmt.Fprintf(out, "")
		fmt.Fprintf(out, "\n")
	}

	fmt.Fprintf(out, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "Options:\n")
	flag.PrintDefaults()
}

func init() {
	flag.IntVar(&port, "port", 5200, "Port for the Proxy Auto Config.")
	flag.IntVar(&port, "p", 5200, "Port for the Proxy Auto Config (shorthand).")
	flag.Var(&tunnelFlag, "link", "Links a port at the host to the given address.")
	flag.Var(&tunnelFlag, "l", "Links a port at the host to the given address (shorthand).")

	flag.Usage = func() {
		usage(os.Stdout, true)
	}
}

func main() {
	flag.Parse()
	fmt.Println(port)
	fmt.Println(tunnelFlag)
	fmt.Printf("%v", flag.Args())
	// server := proxy.NewProxy()
	// listener, _ := net.Listen("tcp", "localhost:0")
	// log.Println(listener.Addr())
	// log.Fatal(server.Serve(listener))
	// Format: sprx -p 5200 -link "some.corporate.web:22 > :8022" -l "localhost:8080 > :8080" alberto@myip.local
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
