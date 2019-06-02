package stockist

/*
Handle Kite Ticket relation operations.

*/

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/beeker1121/goque"
	ms "github.com/mitchellh/mapstructure"
	kiteconnect "github.com/zerodhatech/gokiteconnect"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

//Order details

var (
	ticker *kiteticker.Ticker
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
	err := ticker.Subscribe([]uint32{112129})
	if err != nil {
		fmt.Println("err: ", err)
	}

	ticker.SetMode(kiteticker.ModeFull, []uint32{112129})
}

// Triggered when tick is recevived
func onTick(tick kiteticker.Tick) {
	fmt.Println("Tick received!")

	for i := 0; i < 10; i++ {
		fmt.Printf(" %d st Tick received \n", i)
		time.Sleep(3)
		dticks := dummyTicks()
		EnqueueTick(dticks)
	}

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

//EnqueueTick stores the newly received tick data into a prefix queue
func EnqueueTick(tick *kiteticker.Tick) {

	// tickData := &TickData{
	// 	Open:              tick.OHLC.Open,
	// 	High:              tick.OHLC.High,
	// 	Low:               tick.OHLC.Low,
	// 	Close:             tick.OHLC.Close,
	// 	TotalBuyQuantity:  tick.TotalBuyQuantity,
	// 	TotalSellQuantity: tick.TotalSellQuantity,
	// 	Timestamp:         tick.Timestamp,
	// }

	queue := &Queue{
		Prefix: string(tick.InstrumentToken),
		Data:   tickToMap(tick),
		Path:   "test-queue",
	}

	qc, _ := queue.Create()
	queue.Client = qc
	fmt.Printf("Priting q client here- %v/n", queue)

	err := queue.Insert()
	if err != nil {
		fmt.Printf("Error here")
	}

	fmt.Println("------------------------------------")
	//OrderDetails.QueueInsert(tick)
	// dbClient := InfluxDBClient()
	// fmt.Println(dbClient)
	//InsertTick(dbClient, tickData)

}

//StoreTickInDB pops items from the Stock Queue and adds them to the queue
func StoreTickInDB(prefix string) {
	queue := &Queue{
		Prefix: prefix,
		Path:   "test-queue",
	}

	for {
		qc, err := queue.Create()
		if err != nil {
			fmt.Printf("Error TTTT creating Queue - %v", err)
			continue
		}

		queue.Client = qc
		data, err := queue.Pop()
		if err == goque.ErrEmpty {
			fmt.Printf("Nothing to Pop - %v", err)
			continue
		}

		tick := mapToTick(data)

		db := NewInfluxDB()
		db.InfluxDBClient()

		db.InsertTick(tick)

	}

}
func tickToMap(tick *kiteticker.Tick) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(tick)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func mapToTick(data map[string]interface{}) *kiteticker.Tick {

	var tick *kiteticker.Tick
	err := ms.Decode(data, &tick)
	if err != nil {
		panic(err)
	}

	return tick

}

func dummyTicks() *kiteticker.Tick {

	t := kiteconnect.Time{}
	//t := time.Now()

	ohlc := &kiteticker.OHLC{
		Open:  123,
		Close: 122,
		High:  130,
		Low:   116,
	}
	ticks := &kiteticker.Tick{
		OHLC:              *ohlc,
		InstrumentToken:   1234,
		TotalBuyQuantity:  1000,
		TotalSellQuantity: 800,
		Timestamp:         t,
	}

	return ticks
}
