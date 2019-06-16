package kite

import (
	"math/rand"
	"strconv"
	"time"

	kiteconnect "github.com/sp98/gokiteconnect"
	kiteticker "github.com/sp98/gokiteconnect/ticker"
)

func dummyTicks() *kiteticker.Tick {

	var mode kiteticker.Mode
	mode = "full"

	t := kiteconnect.Time{}
	t.Time = time.Now()

	ohlc := &kiteticker.OHLC{
		Open:  getFloat64(20, 30),
		Close: getFloat64(20, 30),
		High:  getFloat64(20, 30),
		Low:   getFloat64(20, 30),
	}
	ticks := &kiteticker.Tick{
		Mode:               mode,
		InstrumentToken:    112129,
		IsTradable:         true,
		IsIndex:            true,
		LastPrice:          getFloat64(20, 30),
		LastTradedQuantity: getUnit32(7000, 8000),
		TotalBuyQuantity:   getUnit32(7000, 8000),
		TotalSellQuantity:  getUnit32(7000, 8000),
		VolumeTraded:       getUnit32(10000, 15000),
		TotalBuy:           getUnit32(7000, 8000),
		TotalSell:          getUnit32(7000, 8000),
		AverageTradePrice:  getFloat64(20, 30),
		OI:                 getUnit32(20, 30),
		OIDayHigh:          getUnit32(20, 30),
		OIDayLow:           getUnit32(15, 25),

		OHLC: *ohlc,

		Timestamp: t,
	}

	return ticks
}

func random(min int, max int) string {
	rand.Seed(time.Now().UnixNano())
	res := rand.Intn(max-min) + min
	return strconv.Itoa(res)
}

func getFloat64(min int, max int) float64 {
	res := random(min, max)
	f, _ := strconv.ParseFloat(res, 64)
	return f
}

func getUnit32(min int, max int) uint32 {
	// var a uint32
	res := random(min, max)
	u, _ := strconv.ParseUint(res, 10, 32)
	return uint32(u)
}
