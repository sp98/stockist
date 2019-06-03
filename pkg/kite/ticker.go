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

//Order details

var (
	ticker *kiteticker.Ticker
	//Subcriptions to the instruments token 112129
	Subcriptions = []uint32{}
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	fmt.Println("Connected")
	log.Printf("Subcriptions - %+v\n", Subcriptions)
	err := ticker.Subscribe(Subcriptions)
	if err != nil {
		fmt.Println("err: ", err)
	}

	ticker.SetMode(kiteticker.ModeFull, Subcriptions)
}

// Triggered when tick is recevived
func onTick(tick kiteticker.Tick) {
	log.Println("Tick Received frome Kite API")
	StoreTickInDB(&tick)

	// //Run with dummy data when market is closed!
	// for i := 0; i < 1000; i++ {
	// 	time.Sleep(2 * time.Second)
	// 	dticks := dummyTicks()
	// 	StoreTickInDB(dticks)
	// }

	//EnqueueTick(&tick)

}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Printf("Order: %v ", order.OrderID)
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
	log.Printf("Tick received: %+v\n", tick)
	log.Println("---------------------------------")
	db := storage.NewDB("http://localhost:8086", "stockist", "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", strconv.FormatUint(uint64(tick.InstrumentToken), 10), storage.CurrentDate("01022006"))
	//db.Measurement = "test"
	db.InsertTick(tick)

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
