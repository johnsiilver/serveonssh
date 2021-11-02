// Package serveonssh provides a type that allows communicating with services on an SSH endpoint that are listening on a domain socket.
package serveonssh

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// Dialer returns a connection to a unix socket on a remote machine.
type Dialer func() (net.Conn, error)

// Proxy handles forwarding traffic sent on a domain socket over SSH to a remote domain socket.
type Proxy struct {
	sshClient *ssh.Client
	dialer    Dialer
}

// Close closes the underlying SSH client.
func (p Proxy) Close() error {
	return p.sshClient.Close()
}

// Dialer returns the Dialer that opens a connection to the remote Unix socket.
func (p Proxy) Dialer() Dialer {
	return p.dialer
}

// New creates a new Proxy and Dialer. sshEnpoint is the host:port of the remote machine. remoteSocket
// is the path to the Unix socket that the service will be listening to. config is the SSH config needed to dial.
// Proxy is doing the forwarding of our traffic to the remote side. Dialer dials the remote side over SSH.
func New(sshEndpoint, remoteSocket string, config *ssh.ClientConfig) (Proxy, error) {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	client, err := ssh.Dial("tcp", sshEndpoint, config)
	if err != nil {
		return Proxy{}, fmt.Errorf("ssh.Dial failed: %s", err)
	}

	var dial Dialer = func() (net.Conn, error) {
		return client.Dial("unix", remoteSocket)
	}

	return Proxy{sshClient: client, dialer: dial}, nil
}
