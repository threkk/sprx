package forwarding

import (
	"log"
	//	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/crypto/ssh"
)

// Connect Connects to a host using SSH.
func Connect(user, host, password string) *ssh.Client {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return client
}
