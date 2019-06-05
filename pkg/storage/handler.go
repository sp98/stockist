package storage

import (
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

var (
	createDB   = "CREATE DATABASE %s"
	orderQuery = "SELECT * FROM Orders WHERE TradeDate=~/%s/"
	ticksQuery = "SELECT * FROM %s"
	tickCQ     = "CREATE CONTINUOUS QUERY %s ON %s BEGIN %s END"
	tickCQTime = "SELECT FIRST(LastPrice) as Open, MAX(LastPrice) as High, MIN(LastPrice) as Low, LAST(LastPrice) as Close, mean(TotalBuyQuantity) as TotalBuyQuantity, mean(TotalSellQuantity) as TotalSellQuantity INTO %s FROM %s GROUP BY time(%s)"
)

// CreateTickCQ creates a continuous query on Tick Measurement.
func (db DB) CreateTickCQ(tradeInterval string, token string) error {
	cqMeasurement := fmt.Sprintf("%s_%s", db.Measurement, tradeInterval)
	query := fmt.Sprintf(tickCQTime, cqMeasurement, db.Measurement, tradeInterval)
	cquery := fmt.Sprintf(tickCQ, cqMeasurement, db.Name, query)
	q := client.NewQuery(cquery, db.Name, "")
	_, err := db.executeQuery(q)
	if err != nil {
		log.Fatalln("Error creating Tick continuous query - ", err)
		return err
	}
	return nil

}

//GetOrders fetchs the orders to be executed
func (db DB) GetOrders() (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(orderQuery, CurrentDate("2006-01-02")),
		Database: db.Name,
	}
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return nil, err
	}
	return response, nil

}

//GetTicks fetchs the aggregated ticks
func (db DB) GetTicks() (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(ticksQuery, db.Measurement),
		Database: db.Name,
	}

	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return nil, err
	}
	return response, nil

}

//StoreTick saves tick data in influx db
func (db DB) StoreTick(tickData *kiteticker.Tick) error {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.Name,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	tick := *tickData
	fields := map[string]interface{}{
		"LastPrice":         tick.LastPrice,
		"TotalBuyQuantity":  tick.TotalBuyQuantity,
		"TotalSellQuantity": tick.TotalSellQuantity,
	}
	tags := map[string]string{
		"InstrumentToken": fmt.Sprint(tick.InstrumentToken),
	}

	err = db.executePointWrite(bp, db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		log.Fatalln("Error storing tick data to db - ", err)
		return err
	}
	return nil

}
