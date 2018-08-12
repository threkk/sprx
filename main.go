package main

import (
	"context"
	"flag"
	"fmt"
	fw "github.com/threkk/sprx/forwarding"
	"github.com/threkk/sprx/proxy"
	term "golang.org/x/crypto/ssh/terminal"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	// Parse and check the arguments and options.
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Invalid amout of parameters. Expected: 1.\n")
		usage(os.Stderr, false)
		os.Exit(1)
	}

	split := strings.Split(strings.TrimSpace(args[0]), "@")
	if len(split) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid paramenter. Expected: user@host. Got: %v\n", split)
		usage(os.Stderr, false)
		os.Exit(1)
	}

	user := split[0]
	host := split[1]

	// Ask user for the password.
	fmt.Printf("Password:")
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Password could not be read: %s\n", err.Error())
		os.Exit(1)
	}

	// Start SSH connection. If there is no connection, nothing can be done.
	ssh := fw.Connect(user, host, string(pass))
	if ssh == nil {
		fmt.Fprintf(os.Stderr, "Login error.\n")
		os.Exit(1)
	}

	// Start the listeners.
	proxyListener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Proxy listener could not be started: %s\n", err.Error())
		os.Exit(1)
	}

	pacListener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "PAC listener could not be started: %s\n", err.Error())
		os.Exit(1)
	}

	// Get the ports.
	_, proxyPort, err := net.SplitHostPort(proxyListener.Addr().String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid port: %s\n", err.Error())
		os.Exit(1)
	}

	_, pacPort, err := net.SplitHostPort(pacListener.Addr().String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid port: %s\n", err.Error())
		os.Exit(1)
	}

	// Init the proxy server.
	proxyServer := proxy.NewProxy()

	// Use the proxy port to start the PAC handler and the PAC server.
	pacHandler := proxy.NewPacHandler(string(proxyPort))
	pacServer := http.Server{
		Handler: pacHandler,
	}

	// Start all the tunnels.
	// Redirect the local port into the local port in the client.
	proxyHost := net.JoinHostPort("localhost", proxyPort)
	pacHost := net.JoinHostPort("localhost", pacPort)
	proxyTunnel := fw.NewTunnel(proxyHost, proxyHost)
	pacTunnel := fw.NewTunnel(pacHost, pacHost)

	proxyTunnel.Connect(ssh)
	pacTunnel.Connect(ssh)
	for _, t := range tunnelFlag {
		t.Connect(ssh)
	}

	// Start the daemon.
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	// TODO: Timeouts https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// Start proxy server
	go proxyServer.Serve(proxyListener)
	// Start PAC server
	go pacServer.Serve(pacListener)

	// Wait for the signal to stop.
	<-stop

	ssh.Close()
	proxyServer.Shutdown(context.Background())
	pacServer.Shutdown(context.Background())
	os.Exit(0)

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

	// Sources:
	// - https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
	// - https://stackoverflow.com/questions/21417223/simple-ssh-port-forward-in-golang#21655505
	// - https://gist.github.com/codref/473351a24a3ef90162cf10857fac0ff3
}
