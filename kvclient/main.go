package main

import "log"
import "time"
import "math/rand"
import "io"

import "golang.org/x/net/context"
import "google.golang.org/grpc"
import . "github.com/sunwangme/bfgo/api/bfkv"
import . "github.com/sunwangme/bfgo/api/bfgateway"

import "github.com/golang/protobuf/ptypes"

const (
	address = "localhost:50059"
	message = "ping"
)

func Ping(kvclient BfKvServiceClient) {
	resp, err := kvclient.Ping(context.Background(), &BfPingData{Message: message})
	if err != nil {
		log.Fatalf("could not Ping: %v", err)
	}
	log.Printf("recv,%s", resp.Message)
}

func PingStreamC(kvclient BfKvServiceClient) {
	stream, err := kvclient.PingStreamC(context.Background())
	if err != nil {
		log.Fatalf("%v.PingStreamC(_) = _, %v", kvclient, err)
	}

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	pingReq := &BfPingData{Message: message}
	anyReq, err := ptypes.MarshalAny(pingReq)
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
	}
	for i := 0; i < 10; i++ {
		if err := stream.Send(anyReq); err != nil {
			log.Fatalf("Send fail,%v", err)
		}
		log.Printf("send,%s", pingReq.Message)
		s := 500 + rd.Int31n(500) - 1
		time.Sleep(time.Duration(s) * time.Millisecond)
	}

	anyResp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv fail,%v", err)
	}

	pingResp := &BfPingData{}
	if ptypes.Is(anyResp, pingResp) {
		ptypes.UnmarshalAny(anyResp, pingResp)
		log.Printf("recv,%s", pingResp.Message)
	} else {
		log.Fatalf("pingstreamc pong : %v", anyResp)
	}
}

func PingStreamS(kvclient BfKvServiceClient) {
	pingReq := &BfPingData{Message: message}
	anyReq, err := ptypes.MarshalAny(pingReq)
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
	}
	stream, err := kvclient.PingStreamS(context.Background(), anyReq)
	if err != nil {
		log.Fatalf("%v.PingStreamS fail, %v", err)
	}
	for {
		anyResp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("PingStreamS fail, %v", err)
		}

		pingResp := &BfPingData{}
		if ptypes.Is(anyResp, pingResp) {
			ptypes.UnmarshalAny(anyResp, pingResp)
			log.Printf("recv,%s", pingResp.Message)
		} else {
			log.Fatalf("PingStreamS pong : %v", anyResp)
		}
	}
}

func PingStreamCS(kvclient BfKvServiceClient) {
	pingReq := &BfPingData{Message: message}
	anyReq, err := ptypes.MarshalAny(pingReq)
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
	}

	stream, err := kvclient.PingStreamCS(context.Background())
	if err != nil {
		log.Fatalf("PingStreamCS fail, %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			anyResp, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a pong : %v", err)
			}

			pingResp := &BfPingData{}
			if ptypes.Is(anyResp, pingResp) {
				ptypes.UnmarshalAny(anyResp, pingResp)
				log.Printf("recv,%s", pingResp.Message)
			} else {
				log.Fatalf("PingStreamCS pong : %v", anyResp)
			}
		}
	}()

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
	}
	for i := 0; i < 10; i++ {
		if err := stream.Send(anyReq); err != nil {
			log.Fatalf("Send fail,%v", err)
		}
		log.Printf("send,%s", pingReq.Message)
		s := 500 + rd.Int31n(500) - 1
		time.Sleep(time.Duration(s) * time.Millisecond)
	}

	stream.CloseSend()
	<-waitc
}

func main() {
	log.Printf("connect kvserver")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	kvclient := NewBfKvServiceClient(conn)

	log.Printf("===Ping===")
	Ping(kvclient)

	log.Printf("===PingStreamC===")
	PingStreamC(kvclient)

	log.Printf("===PingStreamS===")
	PingStreamS(kvclient)

	log.Printf("===PingStreamCS===")
	PingStreamCS(kvclient)
}
