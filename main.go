package main

import (
	"fmt"

	"github.com/stockist/pkg/stockist"
)

func main() {
	fmt.Println("-- Welcome to Stockist --")

	// Create a DB connection
	db := stockist.NewInfluxDB()
	db.InfluxDBClient()
	fmt.Printf("DB - %+v\n", db)

	// Get Todays Orders from the DB

	//Get Order details
	orderDetails := stockist.NewOrderDetails()
	//Create Kite Connection
	kc, accessToken := stockist.KiteConnect()
	fmt.Printf("AccessToken - %s", accessToken)

	orderDetails.KiteClient = kc

	//Start Kite Ticker:
	go stockist.StartTicker(accessToken)
	//go stockist.StoreTickInDB(orderDetails.InstrumentToken)

	fmt.Scanf("Enter your name here - ")

}
