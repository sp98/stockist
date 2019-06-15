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
	{"YESBANK", "3050241", "NSE", "5m", "20000", "2019-06-14"},
	{"IBULHSGFIN", "7712001", "NSE", "5m", "20000", "2019-06-14"},
	{"ZEEL", "975873", "NSE", "5m", "20000", "2019-06-14"},
	{"BANKBARODA", "1195009", "NSE", "5m", "20000", "2019-06-14"},
	{"SRTRANSFIN", "1102337", "NSE", "5m", "20000", "2019-06-14"},
	{"MOTHERSUMI", "1076225", "NSE", "5m", "20000", "2019-06-14"},
	{"CADILAHC", "2029825", "NSE", "5m", "20000", "2019-06-14"},
	{"ICICIPRULI", "4774913", "NSE", "5m", "20000", "2019-06-14"},
	{"ICICIGI", "5573121", "NSE", "5m", "20000", "2019-06-14"},
	{"ADANIPORTS", "3861249", "NSE", "5m", "20000", "2019-06-14"},
	{"JSWSTEEL", "3001089", "NSE", "5m", "20000", "2019-06-14"},
	{"VEDL", "784129", "NSE", "5m", "20000", "2019-06-14"},
	{"RPOWER", "3906305", "NSE", "5m", "20000", "2019-06-14"},
	{"DISHTV", "3721473", "NSE", "5m", "20000", "2019-06-14"},
	{"SUZLON", "3076609", "NSE", "5m", "20000", "2019-06-14"},
	{"INFRATEL", "7458561", "NSE", "5m", "20000", "2019-06-14"},
	{"DLF", "3771393", "NSE", "5m", "20000", "2019-06-14"},
	{"NTPC", "2977281", "NSE", "5m", "20000", "2019-06-14"},
	{"BPCL", "134657", "NSE", "5m", "20000", "2019-06-14"},
	{"TATASTEEL", "895745", "NSE", "5m", "20000", "2019-06-14"},
	{"INDIGO", "2865921", "NSE", "5m", "20000", "2019-06-14"},
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
