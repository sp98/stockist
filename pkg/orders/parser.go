package orders

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

var (
	orderFields = []string{"Exchange", "InstrumentName", "InstrumentToken", "TradeAmount", "TradeDate", "TradeInterval"}
)

//ResParser parses the InfluxDB result into appropriate Struct
func ResParser(res *client.Response, ord *Order) {
	log.Printf("Parsing result - %+v\n", res)

}
