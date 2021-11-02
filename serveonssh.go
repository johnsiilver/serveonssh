// Package serveonssh provides a type that allows communicating with services on an SSH endpoint that are listening on a domain socket.
package serveonssh

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// Dialer returns a connection to a unix socket.
type Dialer func() (net.Conn, error)

// Forwarder handles forwarding traffic sent on a domain socket over SSH to a remote domain socket.
type Forwarder struct {
	remoteSocket string
	l            net.Listener
	sshClient    *ssh.Client

	mu sync.Mutex
}

// New creates a new Forwarder and Dialer. sshEnpoint is the host:port of the remote machine. remoteSocket
// is the path to the Unix socket that the service will be listening to. config is the SSH config needed to dial.
// Forwarder is doing the forwarding of our traffic to the remote side. Dialer dials the remote side over SSH.
func New(sshEndpoint, remoteSocket string, config *ssh.ClientConfig) (*Forwarder, Dialer, error) {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	client, err := ssh.Dial("tcp", sshEndpoint, config)
	if err != nil {
		return nil, nil, fmt.Errorf("ssh.Dial failed: %s", err)
	}

	f := &Forwarder{
		remoteSocket: remoteSocket,
		sshClient:    client,
	}

	var dial Dialer = func() (net.Conn, error) {
		localSocket := filepath.Join(os.TempDir(), uuid.New().String()+".sock")
		if err := f.localListener(localSocket); err != nil {
			return nil, err
		}
		conn, err := net.DialTimeout("unix", localSocket, 5*time.Second)
		if err != nil {
			log.Println("dial error: ", err)
			return nil, err
		}
		return conn, err
	}

	return f, dial, nil
}

func (f *Forwarder) Close() {
	f.sshClient.Close()
}

func (f *Forwarder) localListener(socket string) error {
	// Setup localListener (type net.Listener)
	l, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
			}
			go f.forward(conn)
		}
	}()
	return nil
}

func (f *Forwarder) forward(localConn net.Conn) {
	sshConn, err := f.sshClient.Dial("unix", f.remoteSocket)
	if err != nil {
		log.Printf("sshConn.Dial() failure: %s", err)
		return
	}

	go func() {
		_, err = io.Copy(sshConn, localConn)
		if err != nil {
			log.Printf("io.Copy failed: %v", err)
		}
	}()

	go func() {
		_, err = io.Copy(localConn, sshConn)
		if err != nil {
			log.Printf("io.Copy failed: %v", err)
		}
	}()
}
