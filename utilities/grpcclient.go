package main

import (
	"flag"
	"io"
	"log"

	"../api"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	serverAddress = flag.String("server", "127.0.0.1:9998", "Grpc server, eg 127.0.0.1:9998")
	streamName    = flag.String("stream", "stream", "Name of stream")
	streamOffset  = flag.Int64("offset", 0, "Offset in stream")
)

func main() {
	flag.Parse()
	log.Printf("Opening stream address=%v, stream=%v", *serverAddress, *streamName)

	connection, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)

	}
	defer connection.Close()

	client := api.NewStreamApiClient(connection)
	arguments := &api.OpenStreamArguments{Id: *streamName, Offset: *streamOffset}
	stream, err := client.OpenStream(context.Background(), arguments)

	if err != nil {
		panic(err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error receiving message %v", err)
		} else {
			log.Printf("Received %v", string(message.Data))
		}
	}
}
