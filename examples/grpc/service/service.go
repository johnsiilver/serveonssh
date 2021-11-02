package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	pb "github.com/johnsiilver/serveonssh/examples/grpc/proto"
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

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterServiceServer(grpcServer, server{})
	grpcServer.Serve(l)
}

type server struct {
	pb.UnimplementedServiceServer
}

func (s server) Say(ctx context.Context, req *pb.Req) (*pb.Resp, error) {
	return &pb.Resp{Say: "World"}, nil
}
