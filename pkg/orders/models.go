package orders

import (
	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

/**
Create all the Models (Struts) here
**/

//Order to be executed
type Order struct {
	InstrumentName  string
	InstrumentToken string
	Exchange        string
	TradeInterval   string
	TradeAmount     string
	TradeDate       string
}

//TickData is the parsed tick data to added to influx DB
type TickData struct {
	Open              float64
	High              float64
	Low               float64
	Close             float64
	TotalBuyQuantity  uint32
	TotalSellQuantity uint32
	Timestamp         kiteconnect.Time
}

// Depth holds the ticker depth data
type Depth struct {
}

var orderData = [][]string{
	{"YESBANK", "136357892", "NSE", "5m", "20000", "2019-06-11"},
	{"IBULHSGFIN", "7712001", "NSE", "5m", "20000", "2019-06-11"},
	{"ZEEL", "975873", "NSE", "5m", "20000", "2019-06-11"},
	{"BANKBARODA", "1195009", "NSE", "5m", "20000", "2019-06-11"},
}

func getOrdersList() *[]Order {
	ordList := []Order{}
	for _, order := range orderData {
		ord := Order{}
		ord.InstrumentName = order[0]
		ord.InstrumentToken = order[1]
		ord.Exchange = order[2]
		ord.TradeInterval = order[3]
		ord.TradeAmount = order[4]
		ord.TradeDate = order[5]
		ordList = append(ordList, ord)
	}

	return &ordList

}
