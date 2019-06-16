package instrument

import (
	"fmt"
	"log"
	"time"

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

//StartProcessing starts the order processing for the day!
func StartProcessing() {
	log.Printf("--- START ---")
	instruments := GetInstrumentList()
	if len(*instruments) == 0 {
		log.Println("Error: No instruments to analyse")
		return
	}
	log.Printf("Instruments - %+v\n", instruments)

	subscriptions := GetSubscriptions()
	if len(subscriptions) == 0 {
		log.Println("Error: No Subscriptions data found")
		return
	}
	log.Printf("Subscriptions - %+v\n", subscriptions)

	//Create Database table and Continous Queries
	err := initDB(*instruments)
	if err != nil {
		log.Printf("Error Initializing DB: %+v", err)
		return
	}

	// Connect to web socket to get tick data.
	for _, subscription := range subscriptions {
		time.Sleep(time.Second * 2)
		go kite.StartTicker(subscription, accessToken)
	}

	// Analyse the aggregated tick results
	for _, instrument := range *instruments {
		cs := &CandleStick{
			Instrument: instrument,
		}
		go cs.startAnalysis()

	}

	//Wait until some time.
	duration, err := closeTrade()
	if err != nil {
		log.Println("Error Parsing Trade Close time")
	}
	log.Printf("Closing trade in - %v", duration)
	time.Sleep(duration)

	log.Printf("--- END ---")
}

func initDB(insturments []Instrument) error {

	//Create DB
	db := storage.NewDB(DBUrl, StockDB, "")
	err := db.CreateDB()
	if err != nil {
		return err
	}

	//Create continuous queries
	for _, instrument := range insturments {
		db := storage.NewDB(DBUrl, StockDB, "")
		db.Measurement = fmt.Sprintf("%s_%s", "ticks", instrument.Token)
		err := db.CreateTickCQ(instrument.Interval, instrument.Token)
		if err != nil {
			log.Printf("Error creating CQ for the isntrument: +%v. Erro: %+v", instrument, err)
			return err
		}
	}
	return nil
}

func closeTrade() (time.Duration, error) {
	closeTime, err := parseTime(layOut, fmt.Sprintf(marketCloseTime, getDate()))
	if err != nil {
		return 0, err
	}
	currentTime := time.Now().Format(tstringFormat)
	parsedCurrentTime, err := parseTime(layOut, currentTime)
	if err != nil {
		return 0, err
	}
	return closeTime.Sub(parsedCurrentTime), nil
}
