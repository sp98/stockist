package instrument

import (
	kiteconnect "github.com/sp98/gokiteconnect"
)

/**
Create all the Models (Struts) here
**/

// Instrument details
type Instrument struct {
	Exchange string
	Symbol   string
	Name     string
	Token    string
	Interval string
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

//GetInstrumentList returns all the instruements to be analysed
func GetInstrumentList() *[]Instrument {
	instrumentList := []Instrument{}
	for _, d := range data {
		inst := Instrument{}
		inst.Name = d[0]
		inst.Symbol = d[1]
		inst.Token = d[2]
		inst.Exchange = d[3]
		inst.Interval = d[4]
		instrumentList = append(instrumentList, inst)
	}

	return &instrumentList

}

// GetSubscriptions gives equally divided subcription list
func GetSubscriptions() [][]uint32 {
	subs := getParsedSubs()
	wbConnections := 3
	var subscriptionList [][]uint32

	chunkSize := (len(subs) + wbConnections - 1) / wbConnections

	for i := 0; i < len(subs); i += chunkSize {
		end := i + chunkSize

		if end > len(subs) {
			end = len(subs)
		}

		subscriptionList = append(subscriptionList, subs[i:end])
	}

	return subscriptionList

}

func getParsedSubs() []uint32 {
	parsedSubs := []uint32{}
	for _, d := range data {
		parsedSubs = append(parsedSubs, getUnit32(d[2]))
	}
	return parsedSubs

}
