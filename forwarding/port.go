// Package forwarding Forwards a series of ports form the
package forwarding

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"net"
)

type Tunnel net.IP
