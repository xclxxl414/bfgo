package bftraderclient

import "log"
import "time"
import "sync/atomic"
import "io"

import "golang.org/x/net/context"
import "google.golang.org/grpc"

import "github.com/golang/protobuf/ptypes"
import "google.golang.org/grpc/metadata"

import . "github.com/sunwangme/bfgo/api/bfgateway"
import . "github.com/sunwangme/bfgo/api/bfdatafeed"
import . "github.com/golang/protobuf/ptypes/any"

//===const===
const (
	addrGateway  = "localhost:50051"
	addrDatafeed = "localhost:50052"
	messagePing  = "ping"
	deadline     = 1
)

//===var===
var (
	spi_      BfTrderClient_ITF = nil
	exitNow_  int32             = 0
	clientId_ string            = "BfTrderClient"

	pingType_         = &BfPingData{}
	accountType_      = &BfAccountData{}
	positionType_     = &BfPositionData{}
	tickType_         = &BfTickData{}
	tradeType_        = &BfTradeData{}
	orderType_        = &BfOrderData{}
	logType_          = &BfLogData{}
	errorType_        = &BfErrorData{}
	notificationType_ = &BfNotificationData{}
)

//===BfTrderClient===
type BfTrderClient struct {
	Gateway      BfGatewayServiceClient
	Datafeed     BfDatafeedServiceClient
	connGateway  *grpc.ClientConn
	connDatafeed *grpc.ClientConn
}

func NewBfTraderClient() *BfTrderClient {
	log.Printf("dail gateway")
	connGateway, err := grpc.Dial(addrGateway, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial gateway fail: %v", err)
	}

	log.Printf("dail datafeed")
	connDatafeed, err := grpc.Dial(addrDatafeed, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial datafeed fail: %v", err)
	}

	gateway := NewBfGatewayServiceClient(connGateway)
	datafeed := NewBfDatafeedServiceClient(connDatafeed)

	return &BfTrderClient{Gateway: gateway, Datafeed: datafeed, connDatafeed: connDatafeed, connGateway: connGateway}
}

//===callback===
func (client *BfTrderClient) OnStart()                                {}
func (client *BfTrderClient) OnNotification(resp *BfNotificationData) {}
func (client *BfTrderClient) OnPing(resp *BfPingData)                 {}
func (client *BfTrderClient) OnTick(resp *BfTickData)                 {}
func (client *BfTrderClient) OnError(resp *BfErrorData)               {}
func (client *BfTrderClient) OnLog(resp *BfLogData)                   {}
func (client *BfTrderClient) OnTrade(resp *BfTradeData)               {}
func (client *BfTrderClient) OnOrder(resp *BfOrderData)               {}
func (client *BfTrderClient) OnPosition(resp *BfPositionData)         {}
func (client *BfTrderClient) OnAccount(resp *BfAccountData)           {}
func (client *BfTrderClient) OnStop()                                 {}

//===gateway api===
func (client *BfTrderClient) SendOrder(req *BfSendOrderReq) (resp *BfSendOrderResp, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	resp, err = client.Gateway.SendOrder(ctx, req)
	return
}
func (client *BfTrderClient) CancleOrder(req *BfCancelOrderReq) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Gateway.CancelOrder(ctx, req)
}
func (client *BfTrderClient) QueryAccount() {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Gateway.QueryAccount(ctx, &BfVoid{})
}
func (client *BfTrderClient) QueryPosition() {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Gateway.QueryPosition(ctx, &BfVoid{})
}
func (client *BfTrderClient) QueryOrders() {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Gateway.QueryOrders(ctx, &BfVoid{})
}
func (client *BfTrderClient) GwGetContract(req *BfGetContractReq) (resps []*BfContractData, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*deadline*time.Second))
	defer cancel()

	resps = make([]*BfContractData, 0)
	stream, err := client.Gateway.GetContract(ctx, req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv fail,%v", err)
			break
		}
		resps = append(resps, resp)
	}

	return resps, nil
}
func (client *BfTrderClient) GwPing(req *BfPingData) (resp *BfPingData, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	resp, err = client.Gateway.Ping(ctx, req)
	return
}

//===datafeed api===
func (client *BfTrderClient) InsertContract(req *BfContractData) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Datafeed.InsertContract(ctx, req)
}
func (client *BfTrderClient) InsertTick(req *BfTickData) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Datafeed.InsertTick(ctx, req)
}
func (client *BfTrderClient) InsertBar(req *BfBarData) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Datafeed.InsertBar(ctx, req)
}
func (client *BfTrderClient) DfGetContract(req *BfGetContractReq) (resps []*BfContractData, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*deadline*time.Second))
	defer cancel()

	resps = make([]*BfContractData, 0)
	stream, err := client.Datafeed.GetContract(ctx, req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv fail,%v", err)
			break
		}
		resps = append(resps, resp)
	}

	return resps, nil
}
func (client *BfTrderClient) GetTick(req *BfGetTickReq) (resps []*BfTickData, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*deadline*time.Second))
	defer cancel()

	resps = make([]*BfTickData, 0)
	stream, err := client.Datafeed.GetTick(ctx, req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv fail,%v", err)
			break
		}
		resps = append(resps, resp)
	}

	return resps, nil
}
func (client *BfTrderClient) GetBar(req *BfGetBarReq) (resps []*BfBarData, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*deadline*time.Second))
	defer cancel()

	resps = make([]*BfBarData, 0)
	stream, err := client.Datafeed.GetBar(ctx, req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv fail,%v", err)
			break
		}
		resps = append(resps, resp)
	}

	return resps, nil
}
func (client *BfTrderClient) DeleteContract(req *BfDeleteContractReq) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Datafeed.DeleteContract(ctx, req)
}
func (client *BfTrderClient) DeleteTick(req *BfDeleteTickReq) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Datafeed.DeleteTick(ctx, req)
}
func (client *BfTrderClient) DeleteBar(req *BfDeleteBarReq) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	client.Datafeed.DeleteBar(ctx, req)

}
func (client *BfTrderClient) DfPing(req *BfPingData) (resp *BfPingData, err error) {
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(deadline*time.Second))
	defer cancel()

	resp, err = client.Datafeed.Ping(ctx, req)
	return
}

//===internal api===
func (client *BfTrderClient) DispatchPush(anyResp *Any) {
	if ptypes.Is(anyResp, tickType_) {
		tickResp := &BfTickData{}
		ptypes.UnmarshalAny(anyResp, tickResp)
		spi_.OnTick(tickResp)
	} else if ptypes.Is(anyResp, pingType_) {
		pingResp := &BfPingData{}
		ptypes.UnmarshalAny(anyResp, pingResp)
		spi_.OnPing(pingResp)
	} else if ptypes.Is(anyResp, accountType_) {
		accountResp := &BfAccountData{}
		ptypes.UnmarshalAny(anyResp, accountResp)
		spi_.OnAccount(accountResp)
	} else if ptypes.Is(anyResp, positionType_) {
		positionResp := &BfPositionData{}
		ptypes.UnmarshalAny(anyResp, positionResp)
		spi_.OnPosition(positionResp)
	} else if ptypes.Is(anyResp, orderType_) {
		orderResp := &BfOrderData{}
		ptypes.UnmarshalAny(anyResp, orderResp)
		spi_.OnOrder(orderResp)
	} else if ptypes.Is(anyResp, tradeType_) {
		tradeResp := &BfTradeData{}
		ptypes.UnmarshalAny(anyResp, tradeResp)
		spi_.OnTrade(tradeResp)
	} else if ptypes.Is(anyResp, logType_) {
		logResp := &BfLogData{}
		ptypes.UnmarshalAny(anyResp, logResp)
		spi_.OnLog(logResp)
	} else if ptypes.Is(anyResp, errorType_) {
		errorResp := &BfErrorData{}
		ptypes.UnmarshalAny(anyResp, errorResp)
		spi_.OnError(errorResp)
	} else if ptypes.Is(anyResp, notificationType_) {
		notificationResp := &BfNotificationData{}
		ptypes.UnmarshalAny(anyResp, notificationResp)
		spi_.OnNotification(notificationResp)
	} else {
		log.Printf("invalid type message,%v", anyResp)
	}
}

func (client *BfTrderClient) ConnectPush(clientId string, tickHandler bool, tradeHandler bool, logHandler bool, symbol string, exchange string) {
	log.Printf("connectPush")

	ctx := context.Background()
	connectPushReq := &BfConnectPushReq{ClientId: clientId, TickHandler: tickHandler, TradeHandler: tradeHandler, LogHandler: logHandler, Symbol: symbol, Exchange: exchange}
	stream, err := client.Gateway.ConnectPush(ctx, connectPushReq)
	if err != nil {
		log.Printf("ConnectPush fail,%v", err)
		return
	}

	for {
		if atomic.CompareAndSwapInt32(&exitNow_, 1, 1) {
			break
		}

		anyResp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv fail,%v", err)
			break
		}

		client.DispatchPush(anyResp)
	}

	log.Printf("connectPush quit")
}

func (client *BfTrderClient) DisconnectPush() {
	log.Printf("disconnectPush")

	ctx := context.Background()
	ctx = metadata.NewContext(ctx, metadata.Pairs("clientid", clientId_))

	_, err := client.Gateway.DisconnectPush(ctx, &BfVoid{})
	if err != nil {
		log.Fatalf("DisconnectPush,%v", err)
	}
}

// detect state by unary rpc with timeout :-(
// https://github.com/grpc/grpc-go/pull/690
func (client *BfTrderClient) DetectServer() bool {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()
	resp, err := client.Gateway.Ping(ctx, &BfPingData{Message: messagePing})
	if err != nil {
		log.Printf("detectServer fail,%v", err)
		return false
	}
	log.Printf("detectServer ok,%s", resp.Message)
	return true
}

func (client *BfTrderClient) FreeConn() {
	client.connDatafeed.Close()
	client.connGateway.Close()
}

//===BfRun===

func BfRun(client BfTrderClient_ITF, clientId string, tickHandler bool, tradeHandler bool, logHandler bool, symbol string, exchange string) {
	log.Printf("start bftraderclient......")
	clientId_ = clientId
	spi_ = client
	firstReady := true

	go monCtrlc(&exitNow_)

	for {
		if atomic.CompareAndSwapInt32(&exitNow_, 1, 1) {
			break
		}

		if client.DetectServer() {
			if firstReady {
				firstReady = false
				client.OnStart()
			}
			client.ConnectPush(clientId, tickHandler, tradeHandler, logHandler, symbol, exchange)
		}

		if atomic.CompareAndSwapInt32(&exitNow_, 1, 1) {
			break
		}
		time.Sleep(time.Duration(5 * time.Second))
	}

	if client.DetectServer() {
		client.OnStop()
		client.DisconnectPush()
	}

	client.FreeConn()

	log.Printf("stop bftraderclient......")
}
