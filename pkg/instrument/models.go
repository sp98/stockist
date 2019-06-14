package instrument

import (
	"os"

	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

/**
Create all the Models (Struts) here
**/

// Instrument details
type Instrument struct {
	Exchange    string
	Symbol      string
	Name        string
	Token       string
	Interval    string
	APIKey      string
	APISecret   string
	AccessToken string
}

// New Instruement
func New() *Instrument {
	instrument := &Instrument{
		Exchange:    os.Getenv("Exchange"),
		Symbol:      os.Getenv("TradingSymbol"),
		Name:        os.Getenv("InstrumentName"),
		Token:       os.Getenv("InstrumentToken"),
		Interval:    os.Getenv("TradeInterval"),
		APIKey:      os.Getenv("APIKey"),
		APISecret:   os.Getenv("APISecret"),
		AccessToken: os.Getenv("AccessToken"),
	}

	return instrument

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
