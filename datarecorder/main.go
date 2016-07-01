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
func (client *DataRecorder) OnNotification(resp *BfNotificationData) {
	log.Printf("OnNotification")
	log.Printf("%v", resp)

	nType := resp.Type
	if nType == BfNotificationType_NOTIFICATION_TRADEWILLBEGIN {
		log.Printf("OnTradeWillBegin")
	} else if nType == BfNotificationType_NOTIFICATION_GOTCONTRACTS {
		log.Printf("OnGotContracts")
		//
		// save contracts
		//
		resps, _ := client.GwGetContract(&BfGetContractReq{Symbol: "*", Exchange: "*"})
		for _, resp := range resps {
			client.InsertContract(resp)
		}
	} else if nType == BfNotificationType_NOTIFICATION_BEGINQUERYORDERS {
	} else if nType == BfNotificationType_NOTIFICATION_BEGINQUERYPOSITION {
	} else if nType == BfNotificationType_NOTIFICATION_ENDQUERYORDERS {
	} else if nType == BfNotificationType_NOTIFICATION_ENDQUERYPOSITION {
	} else {
		log.Printf("invalid notification type")
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
		client.tradeHandler,
		client.logHandler,
		client.symbol,
		client.exchange)
}
