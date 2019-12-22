// sendstdin is an example client that sends logs read from stdin to http://127.0.0.1:9000/.
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/jrockway/logyeet/pkg/logyeetpb"
	"github.com/jrockway/logyeet/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	source := os.Getenv("LOGYEET_SOURCE")
	cc, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := logyeetpb.NewLogyeetClient(cc)

stream:
	for {
		log.Printf("creating new stream %q", source)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.MD{types.SourceHeader: []string{source}})
		conn, err := client.StreamLogs(ctx)
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				log.Printf("writing %d bytes", n)
				req := &logyeetpb.StreamLogsRequest{
					Log: &logyeetpb.Log{
						Log: buf[0:n],
					},
				}
				if err := conn.Send(req); err != nil {
					log.Printf("writing log to stream: %v", err)
					continue stream
				}
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Printf("done reading stdin, exiting")
					_, err := conn.CloseAndRecv()
					if err != nil {
						log.Printf("closing stream: %v", err)
					}
					break stream
				} else {
					log.Printf("read: %v", err)
					os.Exit(1)
				}
			}
		}
	}
	if err := cc.Close(); err != nil {
		log.Printf("closing connection: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
