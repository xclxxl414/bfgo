package main

import "log"
import . "github.com/sunwangme/bfgo/bftraderclient"
import . "github.com/sunwangme/bfgo/api/bfgateway"

//======
type DataRecorder struct {
	*BfTrderClient
	clientId     string
	tickHandler  bool
	tradeHandler bool
	logHandler   bool
	symbol       string
	exchange     string
}

//======
func (client *DataRecorder) OnStart() {
	log.Printf("OnStart")
}
func (client *DataRecorder) OnTradeWillBegin(resp *BfNotificationData) {
	log.Printf("OnTradeWillBegin")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnGotContracts(resp *BfNotificationData) {
	log.Printf("OnGotContracts")
	log.Printf("%v", resp)

	//
	// save contracts
	//
	resps, _ := client.GwGetContract(&BfGetContractReq{Symbol: "*", Exchange: "*"})
	for _, resp := range resps {
		client.InsertContract(resp)
	}
}
func (client *DataRecorder) OnPing(resp *BfPingData) {
	log.Printf("OnPing")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnTick(resp *BfTickData) {
	log.Printf("OnTick")
	log.Printf("%v", resp)

	//
	// save tick
	//
	client.InsertTick(resp)
}
func (client *DataRecorder) OnError(resp *BfErrorData) {
	log.Printf("OnError")
	log.Printf("%v", resp)

}
func (client *DataRecorder) OnLog(resp *BfLogData) {
	log.Printf("OnLog")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnTrade(resp *BfTradeData) {
	log.Printf("OnTrade")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnOrder(resp *BfOrderData) {
	log.Printf("OnOrder")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnPosition(resp *BfPositionData) {
	log.Printf("OnPosition")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnAccount(resp *BfAccountData) {
	log.Printf("OnAccount")
	log.Printf("%v", resp)
}
func (client *DataRecorder) OnStop() {
	log.Printf("OnStop")
}

//======
func main() {
	client := &DataRecorder{
		BfTrderClient: NewBfTraderClient(),
		clientId:      "DataRecorder",
		tickHandler:   true,
		tradeHandler:  false,
		logHandler:    false,
		symbol:        "*",
		exchange:      "*"}

	BfRun(client,
		client.clientId,
		client.tickHandler,
		client.tickHandler,
		client.logHandler,
		client.symbol,
		client.exchange)
}
