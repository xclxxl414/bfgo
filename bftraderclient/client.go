package bftraderclient

import . "github.com/sunwangme/bfgo/api/bfgateway"
import . "github.com/sunwangme/bfgo/api/bfdatafeed"
import . "github.com/golang/protobuf/ptypes/any"

type BfTrderClient_ITF interface {
	//
	// callback, override!!!
	//
	OnStart()
	OnNotification(resp *BfNotificationData)
	OnPing(resp *BfPingData)
	OnTick(resp *BfTickData)
	OnError(resp *BfErrorData)
	OnLog(resp *BfLogData)
	OnTrade(resp *BfTradeData)
	OnOrder(resp *BfOrderData)
	OnPosition(resp *BfPositionData)
	OnAccount(resp *BfAccountData)
	OnStop()

	// gateway api
	SendOrder(req *BfSendOrderReq) (resp *BfSendOrderResp, err error)
	CancleOrder(req *BfCancelOrderReq)
	QueryAccount()
	QueryPosition()
	QueryOrders()
	GwGetContract(req *BfGetContractReq) (resps []*BfContractData, err error)
	GwPing(req *BfPingData) (resp *BfPingData, err error)

	//
	// datafeed api
	//
	InsertContract(req *BfContractData)
	InsertTick(req *BfTickData)
	InsertBar(req *BfBarData)
	DfGetContract(req *BfGetContractReq) (resps []*BfContractData, err error)
	GetTick(req *BfGetTickReq) (resps []*BfTickData, err error)
	GetBar(req *BfGetBarReq) (resps []*BfBarData, err error)
	DeleteContract(req *BfDeleteContractReq)
	DeleteTick(req *BfDeleteTickReq)
	DeleteBar(req *BfDeleteBarReq)
	DfPing(req *BfPingData) (resp *BfPingData, err error)

	//
	// internal
	//
	ConnectPush(clientId string, tickHandler bool, tradeHandler bool, logHandler bool, symbol string, exchange string)
	DispatchPush(resp *Any)
	DisconnectPush()
	DetectServer() bool
	FreeConn()
}
