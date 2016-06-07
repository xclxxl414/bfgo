package main

import "log"

import "golang.org/x/net/context"
import "google.golang.org/grpc"
import . "github.com/sunwangme/bfgo/api/bfgateway"

const (
	address = "localhost:50051"
	message = "bfgo"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewBfGatewayServiceClient(conn)

	r, err := c.Ping(context.Background(), &BfPingData{Message: message})
	if err != nil {
		log.Fatalf("could not Ping: %v", err)
	}
	log.Printf("ctpgateway Pong: %s", r.Message)
}
