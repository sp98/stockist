package storage

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	kiteticker "github.com/zerodhatech/gokiteconnect/ticker"
)

/*
Handle influx db operations in this file
Errors:
1. Database not found:

*/

var (
	createDB   = "CREATE DATABASE %s"
	orderQuery = "SELECT * FROM Orders WHERE TradeDate=~/%s/"
	cQuery     = "CREATE CONTINUOUS QUERY %s ON %s BEGIN %s END"
	cQuery3m   = "SELECT FIRST(LastPrice) as Open, MAX(LastPrice) as High, MIN(LastPrice) as Low, LAST(LastPrice) as Close, mean(TotalBuyQuantity) as TotalBuyQuantity, mean(TotalSellQuantity) as TotalSellQuantity INTO %s FROM %s GROUP BY time(%s)"
)

//DB is the influx db struct
type DB struct {
	Address       string
	Name          string
	Measurement   string
	TradeInterval string
}

//NewDB returns instance of an InfluxDB struct
func NewDB(address string, name string, measurement string) *DB {
	return &DB{
		Address:     address,
		Name:        name,
		Measurement: measurement,
	}

}

//GetClient creates a new Influx DB client
func (db *DB) GetClient() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: db.Address,
	})
	if err != nil {
		log.Fatalln("Error on creating Influx DB client: ", err)
		return nil, err
	}
	return c, nil
}

//CreateDB creates a new database
func (db DB) CreateDB() error {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command: fmt.Sprintf(createDB, db.Name),
	}

	_, err := db.executeQuery(query)
	if err != nil {
		//TODO: Handle dabase already exists error.
		log.Fatalln("Error creating database - ", err)
		return err
	}
	return nil
}

//GetOrders fetchs the orders to be executed
func (db DB) GetOrders() (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(orderQuery, CurrentDate("01-02-2006")),
		Database: db.Name,
	}

	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return nil, err
	}
	return response, nil

}

//InsertTick data in influx db
func (db DB) InsertTick(tickData *kiteticker.Tick) error {
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
		"instrument_token": fmt.Sprint(tick.InstrumentToken),
	}

	pt, err := client.NewPoint(db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		log.Fatalln("Error creating new point - ", err)
		return nil

	}

	bp.AddPoint(pt)

	// Write the batch
	if err := dbClient.Write(bp); err != nil {
		log.Fatalln("Error writing batch - ", err)
		return nil
	}
	return nil

}

// CreateContinuousQuery creates a continuous query on a measurement.
func (db DB) CreateContinuousQuery() error {
	measurement := fmt.Sprintf("%s_%s", db.Measurement, db.TradeInterval)
	query := fmt.Sprintf(cQuery3m, measurement, db.Measurement, db.TradeInterval)
	cquery := fmt.Sprintf(cQuery, measurement, db.Name, query)
	q := client.NewQuery(cquery, db.Name, "")
	_, err := db.executeQuery(q)
	if err != nil {
		log.Fatalln("Error creating continuous error - ", err)
		return err
	}
	return nil

}

func (db DB) executeQuery(query client.Query) (*client.Response, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()

	response, err := dbClient.Query(query)

	if err != nil && response.Error() != nil {
		log.Fatalln("Error executing Query - ", err)
		return nil, err
	}

	return response, nil

}

//CurrentDate returns the date in a specified format
func CurrentDate(format string) string {
	currentTime := time.Now()
	return currentTime.Format(format)
}
