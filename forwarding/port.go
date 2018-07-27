// Package forwarding Forwards a series of ports form the
package forwarding

import (
	"fmt"
	"github.com/threkk/sprx/util"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

const linkRegex = regexp.MustCompile(`([^:\s]+:\d*)\s*>\s*:(\d*)`)

// Tunnel Models the a tunnel between the local machine and the client we are
// connecting.
type Tunnel struct {
	Src string
	Dst string
}

func ParseTunnel(str string) *Tunnel {
	if !linkRegex.Match(str) {
		return nil
	}

	matches := linkRegex.FindStringSubmatch(str)
	src := matches[1]
	dst := fmt.Sprintf("localhost:%s", matches[2])

	tunnel := &Tunnel{Src: src, Dst: dst}
	return tunnel
}

func NewTunnel(dst, src string) *Tunnel {
	tunnel := &Tunnel{Src: src, Dst: dst}
	return tunnel
}

// Connect Given an open and active SSH connection, connects from the local
// machine to the source of the tunnel and using the open SSH connection,
// allocates the requested port at the client and copies the content into it.
func (t *Tunnel) Connect(client *ssh.Client) {
	// Connect to the remote website we will forward. If we cannot connect, we
	// stop the forwarding.
	target, err := net.DialTimeout("tcp", t.Src, 10*time.Second)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Try to allocate the port in the remote client to listen to the content of
	// target. If it is not possible, we skip.
	listener, err := client.Listen("tcp", t.Dst)
	if err != nil {
		log.Fatal(err.Error())
	}

	// If everything goes well, we establish the connection and copy the content
	// of the target into the listener in the remote host.
	res, err := listener.Accept()

	go util.Transfer(res, target)
	go util.Transfer(target, res)
}
