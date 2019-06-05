package storage

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

/*
Handle influx db operations in this file
Errors:
1. Database not found:

*/

//DB is the influx db struct
type DB struct {
	Address     string
	Name        string
	Measurement string
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

func (db DB) executePointWrite(bp client.BatchPoints, measurement string, tags map[string]string, fields map[string]interface{}, t time.Time) error {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	pt, err := client.NewPoint(measurement, tags, fields, t)
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

//CurrentDate returns the date in a specified format
func CurrentDate(format string) string {
	currentTime := time.Now()
	return currentTime.Format(format)
}
