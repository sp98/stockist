package instrument

import (
	"fmt"
	"log"
	"strconv"

	"github.com/stockist1/pkg/kite"
	"github.com/stockist1/pkg/storage"
)

var (
	//DBUrl is the url to connect to the influx DB
	DBUrl = "http://localhost:8086"
	//StockDB is the main database to hold ticks information
	StockDB       = "stockist"
	tradeOpenTime = "%sT03:45:05Z"
)

func initDB(insturment Instrument) error {
	//Create continuous queries
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s", "ticks", insturment.Token)
	err := db.CreateTickCQ(insturment.Interval, insturment.Token)
	if err != nil {
		log.Printf("Error creating CQ for the isntrument: +%v. Erro: %+v", insturment, err)
		return err
	}

	return nil
}

//StartProcessing starts the order processing for the day!
func (instrument Instrument) StartProcessing() {
	log.Println("Start Validating Instrument")
	log.Printf("Instrument - %+v\n", instrument)

	//Create Continous Queries
	err := initDB(instrument)
	if err != nil {
		log.Printf("Instruement valiation failed with Error - %+v", err)
		return
	}
	kite.Subcriptions = append(kite.Subcriptions, getUnit32(instrument.Token))

	// Create connection with Kite
	cs := &CandleStick{
		Instrument: instrument,
	}

	go cs.startAnalysis()

	//Start Kite Ticker
	kite.StartTicker(instrument.APIKey, instrument.AccessToken)

}

func (cs CandleStick) getLowestPrice() float64 {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", cs.Instrument.Token, cs.Instrument.Interval)
	lowest, _ := db.GetLowestLow()
	return lowest
}

func (cs CandleStick) getHighestPrice() float64 {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", cs.Instrument.Token, cs.Instrument.Interval)
	hightest, _ := db.GetMaxHigh()
	return hightest
}

func (cs CandleStick) getOpenPrice() float64 {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", cs.Instrument.Token, cs.Instrument.Interval)
	hightest, _ := db.GetMarketOpenPrice()
	return hightest
}

func getUnit32(str string) uint32 {
	// var a uint32
	u, _ := strconv.ParseUint(str, 10, 32)
	return uint32(u)
}
