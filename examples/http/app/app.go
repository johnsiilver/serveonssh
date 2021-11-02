package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/johnsiilver/serveonssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

var (
	endpoint = flag.String("endpoint", "", "The host:port we are connecting to")
	socket   = flag.String("socket", "", "The Unix socket on the REMOTE side to connect to")
	keyFile  = flag.String("key", "", "The SSH key to use. If not provided, attempts to use the SSH agent.")
	pass     = flag.String("pass", "", "File containing a password to use for SSH. If not provided tries --key and then the SSH agent.")
	user     = flag.String("user", os.Getenv("USER"), "The user to SSH as, set to your logged in user")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	auths, err := getAuthFromFlags()
	if err != nil {
		log.Fatalf("auth failure: %s", err)
	}

	f, dial, err := serveonssh.New(
		*endpoint,
		*socket,
		&ssh.ClientConfig{
			User:            *user,
			Auth:            auths,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Don't do this in real life
		},
	)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return dial()
			},
		},
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := httpc.Get("http://unix" + *socket)
			if err != nil {
				panic(err)
			}

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			if string(b) != "Hello" {
				log.Println("server returned: ", string(b))
			}
			log.Printf("attempt(%d) was successful", i)
		}()
	}

	wg.Wait()
}

func getAuthFromFlags() ([]ssh.AuthMethod, error) {
	auths := []ssh.AuthMethod{}
	if *keyFile != "" {
		a, err := publicKey(*keyFile)
		if err != nil {
			return nil, err
		}
		auths = append(auths, a)
	}
	if *pass != "" {
		b, err := os.ReadFile(*pass)
		if err != nil {
			return nil, fmt.Errorf("pasword file(%s) had error: %s", *pass, err)
		}
		auths = append(auths, ssh.Password(strings.TrimSpace(string(b))))
	}
	if a, err := agentAuth(); err == nil {
		auths = append(auths, a)
	}
	return auths, nil
}

func agentAuth() (ssh.AuthMethod, error) {
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, err
	}

	client := agent.NewClient(conn)
	return ssh.PublicKeysCallback(client.Signers), nil
}

func publicKey(privateKeyFile string) (ssh.AuthMethod, error) {
	k, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(k)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
