package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	socketAddr = flag.String("socket", "", "The unix socket to open")
)

func main() {
	flag.Parse()

	switch *socketAddr {
	case "/":
		log.Fatalf("can't have a socket at root")
	case "":
		log.Fatalf("socket must be set")
	}

	if err := os.MkdirAll(filepath.Dir(*socketAddr), 0700); err != nil {
		log.Printf("could not create socket dir path: %w", err)
	}
	if err := os.Remove(*socketAddr); err != nil {
		log.Printf("could not remove old socket file: ", err)
	}

	l, err := net.Listen("unix", *socketAddr)
	if err != nil {
		log.Printf("could not connect to socket: %w", err)
	}

	h := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello!")
		},
	)

	serv := &http.Server{Handler: h}

	log.Fatal(serv.Serve(l))
}
