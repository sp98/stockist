package kite

/*
Handle Kite Ticket relation operations.

*/

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stockist/pkg/storage"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

//TDepth is the depth of the ticker.
type TDepth kiteticker.Depth

//Order details

var (
	ticker *kiteticker.Ticker
	//Subcriptions to the instruments token 112129
	Subcriptions   = []uint32{}
	marketOpenTime = "%s 9:00:00"
	DBUrl          = "http://influxdb:8086"
	//StockDB is the main database to hold ticks information
	StockDB = "stockist"
)

// Triggered when any error is raised
func onError(err error) {
	log.Println("Error in Kite Trade API: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	log.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	log.Println("Connected with Kite Trading API")
	log.Printf("Subcriptions - %+v\n", Subcriptions)
	err := ticker.Subscribe(Subcriptions)
	if err != nil {
		fmt.Println("err: ", err)
	}

	ticker.SetMode(kiteticker.ModeFull, Subcriptions)
}

// Triggered when tick is recevived
func onTick(tick kiteticker.Tick) {
	//log.Println("Tick Received frome Kite API")
	//StoreTickInDB(&tick)

	//Run with dummy data when market is closed!
	for i := 0; i < 1000; i++ {
		time.Sleep(2 * time.Second)
		dticks := dummyTicks()
		StoreTickInDB(dticks)
	}

}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	log.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	log.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	log.Printf("Order: %+v ", order)
	if order.Status == "COMPLETE" && order.TransactionType == "BUY" && isInstrumentSubscribed(order.InstrumentToken) {
		updateTradeInDB("BOUGHT", strconv.FormatUint(uint64(order.InstrumentToken), 10))
	} else if order.Status == "COMPLETE" && order.TransactionType == "SELL" && isInstrumentSubscribed(order.InstrumentToken) {
		updateTradeInDB("SOLD", strconv.FormatUint(uint64(order.InstrumentToken), 10))
	} else if order.Status == "REJECTED" && isInstrumentSubscribed(order.InstrumentToken) {
		log.Printf("Last Order Got Rejected")
	}
}

//StartTicker starts the websocket to receive kite ticker data
func StartTicker(apiKey, accestoken string) {
	// Create a new Kite connect instance

	// Create new Kite ticker instance
	ticker = kiteticker.New(apiKey, accestoken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)
	ticker.OnOrderUpdate(onOrderUpdate)

	// Start the connection
	ticker.Serve()

}

//StoreTickInDB stors the tick in influx db
func StoreTickInDB(tick *kiteticker.Tick) {
	log.Printf("Tick received: %+v\n", tick)
	log.Println("---------------------------------")
	db := storage.NewDB(DBUrl, StockDB, "")
	if isBeforeMarketOpen() {
		db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", strconv.FormatUint(uint64(tick.InstrumentToken), 10), "5m")
		log.Println("Storing Last Day's OHLC ")
		db.StorePreviousDayOHLC(tick)
	} else {
		db.Measurement = fmt.Sprintf("%s_%s", "ticks", strconv.FormatUint(uint64(tick.InstrumentToken), 10))
		//TDepth = &tick.Depth{}
		db.StoreTick(tick)
	}

}

// Check if the current time is before the market open time
func isBeforeMarketOpen() bool {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	mot := fmt.Sprintf(marketOpenTime, time.Now().Format("2006-01-02"))
	//t, err := time.Parse("2006-01-02 15:04:05", mot)
	t, err := time.ParseInLocation("2006-01-02 15:04:05", mot, loc)

	if err != nil {
		fmt.Println(err)
	}

	return time.Now().Before(t)

}

func updateTradeInDB(option, instToken string) {
	db := storage.NewDB(DBUrl, StockDB, "trade")
	db.InsertTrade(instToken, option)

}

func isInstrumentSubscribed(instToken uint32) bool {
	for _, subscribedInst := range Subcriptions {
		if instToken == subscribedInst {
			//log.Println("Token Matched")
			return true
		}

	}
	//log.Println("Token didn't match!")
	return false
}
