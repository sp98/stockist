package orders

import (
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/stockist/pkg/kite"
	"github.com/stockist/pkg/storage"
)

var (
	//DBUrl is the url to connect to the influx DB
	DBUrl = "http://localhost:8086"
	//StockDB is the main database to hold ticks information
	StockDB       = "stockist"
	tradeOpenTime = "%sT03:45:05Z"
)

func initDB(order Order) error {
	//Create continuous queries
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s", "ticks", order.InstrumentToken)
	err := db.CreateTickCQ(order.TradeInterval, order.InstrumentToken)
	if err != nil {
		return err
	}

	return nil
}

//StartProcessing starts the order processing for the day!
func StartProcessing() {
	log.Println("Start Processing Order")
	// Check for valid orders to be processed today
	//orders := GetOrders()
	orders := getOrdersList()
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

	for _, order := range *orders {
		trade := &Trade{
			Order: order,
		}

		go trade.startAnalysis()
	}

	//Start Kite Ticker
	kite.StartTicker(accessToken)

}

//GetOrders from the Orders Measurement in the databse
func GetOrders() *[]Order {
	db := storage.NewDB(DBUrl, StockDB, "")
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

func (trade Trade) getLowestPrice() float64 {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", trade.Order.InstrumentToken, trade.Order.TradeInterval)
	lowest, _ := db.GetLowestLow()
	return lowest
}

func (trade Trade) getHighestPrice() float64 {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", trade.Order.InstrumentToken, trade.Order.TradeInterval)
	hightest, _ := db.GetMaxHigh()
	return hightest
}

func (trade Trade) getOpenPrice() float64 {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", trade.Order.InstrumentToken, trade.Order.TradeInterval)
	hightest, _ := db.GetMarketOpenPrice()
	return hightest
}

func getUnit32(str string) uint32 {
	// var a uint32
	u, _ := strconv.ParseUint(str, 10, 32)
	return uint32(u)
}
