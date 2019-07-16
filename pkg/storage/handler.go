package storage

import (
	"fmt"
	"log"
	"strconv"
	"time"

	client "github.com/orourkedd/influxdb1-client/client"
	kiteticker "github.com/sp98/gokiteconnect/ticker"
)

var (
	createDB   = "CREATE DATABASE %s"
	orderQuery = "SELECT * FROM Orders WHERE TradeDate=~/%s/"
	tradeQuery = `select * from trade where InstrumentToken='%s' ORDER BY time DESC limit 1`

	firstCandleStickQuery = `select * from %s limit 1`
	maxHighQuery          = "SELECT max(High) as Highest from %s"
	minLowQuery           = "SELECT min(Low) as Lowest from %s"
	ticksQuery            = "SELECT * FROM %s ORDER BY time DESC"
	ohlcQuery             = "SELECT * FROM %s ORDER BY time DESC limit 1"
	tickCQ                = "CREATE CONTINUOUS QUERY %s ON %s BEGIN %s END"
	tickCQTime            = "SELECT FIRST(LastPrice) as Open, MAX(LastPrice) as High, MIN(LastPrice) as Low, LAST(LastPrice) as Close, last(AverageTradePrice) as AverageTradePrice INTO %s FROM %s GROUP BY time(%s)"
)

// CreateTickCQ creates a continuous query on Tick Measurement.
func (db DB) CreateTickCQ(tradeInterval, token string) error {
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

//GetOHLC returns Open High Low and Close value for a Stock
func (db DB) GetOHLC() (float64, float64, float64, float64, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(ohlcQuery, db.Measurement),
		Database: db.Name,
	}
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return 0, 0, 0, 0, err
	}

	//log.Printf("Res - %+v", response)
	Open := response.Results[0].Series[0].Values[0][7]
	Openf, _ := strconv.ParseFloat(fmt.Sprintf("%v", Open), 64)
	High := response.Results[0].Series[0].Values[0][4]
	Highf, _ := strconv.ParseFloat(fmt.Sprintf("%v", High), 64)
	Low := response.Results[0].Series[0].Values[0][6]
	Lowf, _ := strconv.ParseFloat(fmt.Sprintf("%v", Low), 64)
	Close := response.Results[0].Series[0].Values[0][3]
	Closef, _ := strconv.ParseFloat(fmt.Sprintf("%v", Close), 64)

	//log.Println("OHLC - ", Open, High, Low, Close)
	return Openf, Highf, Lowf, Closef, nil

}

//GetTradeQuantity returns the Buy and sell quanity
func (db DB) GetTradeQuantity() (float64, float64, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		Command:  fmt.Sprintf(ohlcQuery, db.Measurement),
		Database: db.Name,
	}
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		return 0, 0, err
	}
	bq := response.Results[0].Series[0].Values[0][2]
	bqf, _ := strconv.ParseFloat(fmt.Sprintf("%v", bq), 64)
	sq := response.Results[0].Series[0].Values[0][8]
	sqf, _ := strconv.ParseFloat(fmt.Sprintf("%v", sq), 64)

	//log.Println("OHLC - ", Open, High, Low, Close)
	return bqf, sqf, nil

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
		"AverageTradePrice": tick.AverageTradePrice,
		"Open":              tick.OHLC.Open,
		"High":              tick.OHLC.High,
		"Low":               tick.OHLC.Low,
		"Close":             tick.OHLC.Close,
		"BuyQuantity":       tick.TotalBuyQuantity,
		"SellQuantity":      tick.TotalSellQuantity,
	}
	tags := map[string]string{
		// "InstrumentToken": fmt.Sprint(tick.InstrumentToken),
	}

	err = db.executePointWrite(bp, db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		log.Fatalln("Error storing tick data to db - ", err)
		return err
	}
	return nil

}

//StorePreviousDayOHLC stores ohlc for previous day
func (db DB) StorePreviousDayOHLC(tickData *kiteticker.Tick) error {
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
		"Open":              tick.OHLC.Open,
		"High":              tick.OHLC.High,
		"Low":               tick.OHLC.Low,
		"Close":             tick.OHLC.Close,
		"AverageTradePrice": tick.AverageTradePrice,
	}
	tags := map[string]string{}

	err = db.executePointWrite(bp, db.Measurement, tags, fields, tickData.Timestamp.Time)
	if err != nil {
		log.Fatalln("Error storing tick data to db - ", err)
		return err
	}
	return nil

}

// GetMaxHigh gives the maximum High so far
func (db DB) GetMaxHigh() (float64, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		//Command:  fmt.Sprintf(maxHighQuery, "ticks_3699201_1m"),
		Command:  fmt.Sprintf(maxHighQuery, db.Measurement),
		Database: db.Name,
	}
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		// return nil, err
	}
	log.Printf("res %+v", response)
	// return response, nil
	if len(response.Results) == 0 {
		return 0, fmt.Errorf("Error finding max High from the aggregared query")
	}
	log.Printf("res %+v", response)
	highestHigh := response.Results[0].Series[0].Values[0][1]
	hightestHighf, _ := strconv.ParseFloat(fmt.Sprintf("%v", highestHigh), 64)
	return hightestHighf, nil

}

//GetLowestLow fetches the lowest value after 9:15 am
func (db DB) GetLowestLow() (float64, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		//Command:  fmt.Sprintf(lowestLowQuery, "ticks_3699201_1m"),
		Command:  fmt.Sprintf(minLowQuery, db.Measurement),
		Database: db.Name,
	}
	//log.Println("Query - ", query.Command)
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		// return nil, err
	}
	// return response, nil
	if len(response.Results) == 0 {
		return 0, fmt.Errorf("Error finding max High from the aggregared query")
	}

	//log.Printf("Response - %+v", response.Results)
	lowestLow := response.Results[0].Series[0].Values[0][1]
	lowestLowf, _ := strconv.ParseFloat(fmt.Sprintf("%v", lowestLow), 64)
	fmt.Println(lowestLowf)
	return lowestLowf, nil
}

//GetMarketOpenPrice gets the opening price
func (db DB) GetMarketOpenPrice() (float64, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		//Command:  fmt.Sprintf(firstCandleStickQuery, db.Measurement, t),
		Command:  fmt.Sprintf(firstCandleStickQuery, "ticks_3699201_1m"),
		Database: db.Name,
	}

	//log.Println("Query - ", query.Command)
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		// return nil, err
	}

	//fmt.Printf("res - %+v", response)
	// return response, nil
	if len(response.Results) == 0 {
		return 0, fmt.Errorf("Error finding max High from the aggregared query")
	}

	open := response.Results[0].Series[0].Values[0][5]
	openf, _ := strconv.ParseFloat(fmt.Sprintf("%v", open), 64)
	fmt.Println(openf)
	return 0, nil
}

// InsertTrade inserts data when a stock needs to be bought
func (db DB) InsertTrade(tokenID string, trade string) error {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db.Name,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	fields := map[string]interface{}{
		"trade": trade,
	}
	tags := map[string]string{
		"InstrumentToken": tokenID,
	}

	err = db.executePointWrite(bp, db.Measurement, tags, fields, time.Now())
	if err != nil {
		log.Fatalln("Error storing trade data to db - ", err)
		return err
	}
	return nil

}

//GetLastTrade gives the last trade done on an Instrument
func (db DB) GetLastTrade(tokenID string) (string, error) {
	dbClient, _ := db.GetClient()
	defer dbClient.Close()
	query := client.Query{
		//Command:  fmt.Sprintf(lowestLowQuery, "ticks_3699201_1m"),
		Command:  fmt.Sprintf(tradeQuery, tokenID),
		Database: db.Name,
	}
	//log.Println("Query - ", query.Command)
	response, err := db.executeQuery(query)
	if err != nil {
		log.Fatalln("Error getting orders - ", err)
		// return nil, err
	}
	// return response, nil
	if len(response.Results[0].Series) == 0 {
		return "", fmt.Errorf("No data found")
	}

	log.Printf("Response - %+v", response.Results)
	lastTrade := response.Results[0].Series[0].Values[0][2]
	slastTrade := fmt.Sprintf("%v", lastTrade)
	fmt.Println(slastTrade)
	return slastTrade, nil

}
