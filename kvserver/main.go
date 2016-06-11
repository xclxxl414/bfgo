package main

import "log"

import "time"
import "math/rand"
import "io"
import "net"

import "golang.org/x/net/context"
import "google.golang.org/grpc"
import . "github.com/sunwangme/bfgo/api/bfkv"
import . "github.com/sunwangme/bfgo/api/bfgateway"

import "github.com/golang/protobuf/ptypes"
import . "github.com/golang/protobuf/ptypes/any"

const (
	address = "localhost:50059"
	message = "pong"
)

//======
type KvServer struct {
}

//======
func (kvserver *KvServer) Ping(ctx context.Context, req *BfPingData) (*BfPingData, error) {
	log.Printf("===Ping===")
	log.Printf("recv,%s", req.Message)
	return &BfPingData{Message: message}, nil
}

func (kvserver *KvServer) PingStreamC(stream BfKvService_PingStreamCServer) error {
	log.Printf("===PingStreamC===")
	pingResp := &BfPingData{Message: message}
	anyResp, err := ptypes.MarshalAny(pingResp)
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
		return err
	}

	for {
		anyReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(anyResp)
		}
		if err != nil {
			log.Fatalf("PingStreamC fail,%v", err)
			return err
		}

		pingReq := &BfPingData{}
		if ptypes.Is(anyReq, pingReq) {
			ptypes.UnmarshalAny(anyReq, pingReq)
			log.Printf("recv,%s", pingReq.Message)
		} else {
			log.Fatalf("PingStreamC,%v", anyReq)
			return err
		}
	}
}
func (kvserver *KvServer) PingStreamS(anyReq *Any, stream BfKvService_PingStreamSServer) error {
	log.Printf("===PingStreamS===")

	pingReq := &BfPingData{}
	if ptypes.Is(anyReq, pingReq) {
		ptypes.UnmarshalAny(anyReq, pingReq)
		log.Printf("recv,%s", pingReq.Message)
	} else {
		log.Fatalf("PingStreamS,%v", anyReq)
		return nil
	}

	pingResp := &BfPingData{Message: message}
	anyResp, err := ptypes.MarshalAny(pingResp)
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
		return err
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		log.Printf("send,%d", i)
		if err := stream.Send(anyResp); err != nil {
			return err
		}
		s := 500 + rd.Int31n(500) - 1
		time.Sleep(time.Duration(s) * time.Millisecond)
	}

	return nil
}

func (kvserver *KvServer) PingStreamCS(stream BfKvService_PingStreamCSServer) error {
	log.Printf("===PingStreamCS===")
	pingResp := &BfPingData{Message: message}
	anyResp, err := ptypes.MarshalAny(pingResp)
	if err != nil {
		log.Fatalf("MarshalAny fail,%v", err)
		return err
	}

	for {
		anyReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		pingReq := &BfPingData{}
		if ptypes.Is(anyReq, pingReq) {
			ptypes.UnmarshalAny(anyReq, pingReq)
			log.Printf("recv,%s", pingReq.Message)
		} else {
			log.Fatalf("PingStreamS,%v", anyReq)
			return err
		}

		if err := stream.Send(anyResp); err != nil {
			return err
		}
	}
}

func (kvserver *KvServer) SetKv(context.Context, *BfKvData) (*BfVoid, error) {
	return &BfVoid{}, nil

}

func (kvserver *KvServer) GetKv(context.Context, *BfKvData) (*BfKvData, error) {
	return &BfKvData{}, nil
}

//======
func newKvServer() *KvServer {
	s := new(KvServer)
	return s
}

//======
func main() {
	log.Printf("kvserver listening on %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterBfKvServiceServer(grpcServer, newKvServer())
	grpcServer.Serve(lis)
}
