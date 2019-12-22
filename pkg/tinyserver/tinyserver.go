// Package tinyserver implements a tiny server that receives logs from a logyeet client.
package tinyserver

import (
	"errors"
	"io"
	"log"

	"github.com/jrockway/logyeet/pkg/logyeetpb"
	"github.com/jrockway/logyeet/pkg/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Print is a logyeet server that prints each log message when it arrives.
type Print struct{}

// StreamLogs implements logyeetpb.LogyeetServer.
func (Print) StreamLogs(stream logyeetpb.Logyeet_StreamLogsServer) error {
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	source := "unknown"
	if ok {
		for _, s := range md.Get(types.SourceHeader) {
			source = s
		}
	}
	log.Printf("new stream: %q", source)
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("[%s] done (ok)", source)
				return stream.SendAndClose(&logyeetpb.StreamLogsReply{})
			}
			log.Printf("[%s] done (%v)", source, err)
			return status.Error(codes.Unknown, err.Error())
		}
		if text := msg.GetLog().GetLog(); len(text) > 0 {
			log.Printf("[%s] %s", source, text)
		}
	}
}
