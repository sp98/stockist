package stockist

import (
	"github.com/beeker1121/goque"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

/**
Create all the Models (Struts) here
**/

//OrderDetails to be executed
type OrderDetails struct {
	KiteClient      *kiteconnect.Client
	DB              *InfluxDB
	Queue           *goque.PrefixQueue
	InstrumentName  string
	InstrumentToken string
	Exchange        string
	TradeInterval   string
	TradeAmount     string
	TradeDate       string
	TickData        *kiteticker.Tick
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

//NewInfluxDB returns instance of an InfluxDB struct
func NewInfluxDB() *InfluxDB {

	return &InfluxDB{
		Address: "http://localhost:8086",
		Name:    "Stockist",
	}

}

//NewOrderDetails gets details about orders to be executed today
func NewOrderDetails() *OrderDetails {
	return &OrderDetails{}
}
