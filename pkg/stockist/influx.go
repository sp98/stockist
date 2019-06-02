package stockist

import (
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

/*
Handle influx db operations in this file
Errors:
1. Database not found:

*/

var cQuery = "CREATE CONTINUOUS QUERY %s ON %s BEGIN %s END"

var cQuery3m = "SELECT mean(Open) as Open, mean(High) as High, mean(Low) as Low, mean(Close) as Close, mean(TotalBuyQuantity) as TotalBuyQuantity, mean(TotalSellQuantity) as TotalBuyQuantity INTO %s FROM %s GROUP BY time(%s)"

//InfluxDB is the influx db struct
type InfluxDB struct {
	Address     string
	Username    string
	Password    string
	Client      *client.Client
	Name        string
	Measurement string
}

//NewDB returns instance of an InfluxDB struct
func NewDB() *InfluxDB {

	return &InfluxDB{
		Address: "http://localhost:8086",
		Name:    "stockist",
	}

}

//InfluxDBClient creates a new Influx DB client
func (db *InfluxDB) InfluxDBClient() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: db.Address,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	db.Client = &c
}

//CreateDB creates a new database
func CreateDB() {

}

//CreateMeasurement creates a new measurement
func CreateMeasurement() {

}

//TodaysOrders fetchs orders to be executed today
func (db InfluxDB) TodaysOrders() {

}

//InsertTick data in influx db
func (db InfluxDB) InsertTick(tickData *kiteticker.Tick) {
	// Create a new point batch

	c := *db.Client
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.Name,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	tick := *tickData
	fields := map[string]interface{}{
		"Open":              tick.OHLC.Open,
		"High":              tick.OHLC.High,
		"Low":               tick.OHLC.Low,
		"Close":             tick.OHLC.Close,
		"TotalBuyQuantity":  tick.TotalBuyQuantity,
		"TotalSellQuantity": tick.TotalSellQuantity,
	}
	tags := map[string]string{
		"instrument_token": fmt.Sprint(tick.InstrumentToken),
	}

	pt, err := client.NewPoint(db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

}

// CreateContinuousQuery creates a continueous query on a measure.
func (db InfluxDB) CreateContinuousQuery() {
	c := *db.Client
	defer c.Close()
	cqMeaurement := fmt.Sprintf("%s_%s", db.Measurement, "2m")
	cq := fmt.Sprintf(cQuery3m, cqMeaurement, db.Measurement, "2m")
	finalQuery := fmt.Sprintf(cQuery, cqMeaurement, db.Name, cq)
	fmt.Printf("CQ Final Query - %s\n", finalQuery)
	q := client.NewQuery(finalQuery, db.Name, "")
	response, err := c.Query(q)

	if err != nil && response.Error() != nil {
		fmt.Printf("Error ins creating CQ - %+v\n", err)
	}
	fmt.Printf("CQ create response - %+v\n", response.Results)

}
