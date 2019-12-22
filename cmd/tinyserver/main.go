// tinyserver is a logyeet server that listens on localhost:9000 and prints the logs it receives.
package main

import (
	"log"
	"net"

	"github.com/jrockway/logyeet/pkg/logyeetpb"
	"github.com/jrockway/logyeet/pkg/tinyserver"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	logyeetpb.RegisterLogyeetServer(s, tinyserver.Print{})
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
