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
	Subcriptions = []uint32{}
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
	StoreTickInDB(&tick)

	//Run with dummy data when market is closed!
	// for i := 0; i < 1000; i++ {
	// 	time.Sleep(2 * time.Second)
	// 	dticks := dummyTicks()
	// 	StoreTickInDB(dticks)
	// }

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
	log.Printf("Order: %v ", order.OrderID)
}

//StartTicker starts the websocket to receive kite ticker data
func StartTicker(accestoken string) {
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
	//log.Printf("Tick received: %+v\n", tick)
	//log.Println("---------------------------------")
	db := storage.NewDB("http://localhost:8086", "stockist", "")
	db.Measurement = fmt.Sprintf("%s_%s", "ticks", strconv.FormatUint(uint64(tick.InstrumentToken), 10))
	//TDepth = &tick.Depth{}
	db.StoreTick(tick)

}

// func tickToMap(tick *kiteticker.Tick) map[string]interface{} {
// 	var inInterface map[string]interface{}
// 	inrec, _ := json.Marshal(tick)
// 	json.Unmarshal(inrec, &inInterface)
// 	return inInterface
// }

// func mapToTick(data map[string]interface{}) *kiteticker.Tick {

// 	var tick *kiteticker.Tick
// 	err := ms.Decode(data, &tick)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return tick

// }
