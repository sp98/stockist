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
		Exchange:    os.Getenv("EXCHANGE"),
		Symbol:      os.Getenv("TRADINGSYMBOL"),
		Name:        os.Getenv("INSTRUMENTNAME"),
		Token:       os.Getenv("INSTRUMENTTOKEN"),
		Interval:    os.Getenv("TRADEINTERVAL"),
		APIKey:      os.Getenv("APIKEY"),
		APISecret:   os.Getenv("APISECRET"),
		AccessToken: os.Getenv("ACCESSTOKEN"),
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
