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
