package orders

import (
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/stockist/pkg/kite"
	"github.com/stockist/pkg/storage"
)

func initDB(order Order) error {
	//Create continuous queries
	db := storage.NewDB("http://localhost:8086", "stockist", "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", order.InstrumentToken, storage.CurrentDate("01022006"))
	db.TradeInterval = order.TradeInterval
	err := db.CreateContinuousQuery()
	if err != nil {
		return err
	}

	return nil
}

//Intialize a new order
func Intialize() {
	GetOrders()

}

//StartProcessing of a new order
func StartProcessing() {
	log.Println("Start Processing Order")
	// Check for valid orders to be processed today
	orders := GetOrders()
	if orders == nil {
		log.Println("No Orders to execute today")
		return
	}
	log.Printf("Orders to be executed - %+v\n", orders)

	//Create Continous Queries
	for _, order := range *orders {
		initDB(order)
		kite.Subcriptions = append(kite.Subcriptions, getUnit32(order.InstrumentToken))
	}

	// Create connection with Kite
	_, accessToken := kite.Connect()
	//Start Kite Ticker:
	kite.StartTicker(accessToken)

}

//GetOrders from the Orders Measurement in the databse
func GetOrders() *[]Order {
	db := storage.NewDB("http://localhost:8086", "stockist", "")
	var orderRespsonse *client.Response
	orderRespsonse, _ = db.GetOrders()

	if len(orderRespsonse.Results[0].Series) == 0 {
		return nil
	}
	var ordList []Order
	for _, results := range orderRespsonse.Results {
		for _, rows := range results.Series {
			ord := &Order{}
			for _, row := range rows.Values {

				ord.Exchange = fmt.Sprintf("%v", row[1])
				ord.InstrumentName = fmt.Sprintf("%v", row[2])
				ord.InstrumentToken = fmt.Sprintf("%v", row[3])
				ord.TradeAmount = fmt.Sprintf("%v", row[4])
				ord.TradeDate = fmt.Sprintf("%v", row[5])
				ord.TradeInterval = fmt.Sprintf("%v", row[6])
				ordList = append(ordList, *ord)

			}

		}

	}

	return &ordList

}

func getUnit32(str string) uint32 {
	// var a uint32
	u, _ := strconv.ParseUint(str, 10, 32)
	return uint32(u)
}