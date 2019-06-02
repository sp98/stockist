package stockist

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

/*
Handle influx db operations in this file
*/

//InfluxDB is the influx db struct
type InfluxDB struct {
	Address     string
	Username    string
	Password    string
	Client      *client.Client
	Name        string
	Measurement string
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
		"instrument_token": "test",
	}

	pt, err := client.NewPoint(db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	c := *db.Client
	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

}
