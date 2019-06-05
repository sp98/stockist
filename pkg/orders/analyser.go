package orders

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stockist/pkg/storage"
)

var (
	marketCloseTime = "%s 23:59:00"
	tstringFormat   = "2006-01-02 15:04:05"
	layOut          = "2006-01-02 15:04:05"
)

//Trade holds the currently trading order details
type Trade struct {
	Order             Order
	Details           []TradeDetails
	IsBullish         bool
	CurrentTrend      string
	CurrentTrendCount int
}

//TradeDetails is the aggregate of the trade
type TradeDetails struct {
	Close             float64
	High              float64
	Low               float64
	Open              float64
	TotalBuyQuantity  float64
	TotalSellQuantity float64
}

func (trade *Trade) startAnalysis() error {
	log.Println("--- Begin Analysis ----")

	interval := getInterval(trade.Order.TradeInterval)
	if interval == 0 {
		log.Fatal("invalid order interval ")
		//return error here
	}

	mct, err := parseTime(layOut, fmt.Sprintf(marketCloseTime, trade.Order.TradeDate))
	if err != nil {
		return err
	}

	if stop, _ := trade.stopAnalysis(mct, false); stop {
		log.Printf("Can't start analysis. Already past market closing time %+v", mct)
		return nil
	}

	wt := waitBeforeAnalysis(interval)
	if wt > 0 {
		log.Printf("Wait for %.2f minutes before starting", float64(wt)/60)
		time.Sleep(time.Second * time.Duration(wt))
	}

	t := time.NewTicker(time.Minute * time.Duration(interval))

	log.Printf("Analysis Start Time: %+v ::: Analysis Stop Time: %+v", time.Now(), fmt.Sprintf(marketCloseTime, trade.Order.TradeDate))

	for alive := true; alive; {
		res, _ := trade.stopAnalysis(mct, true)
		if res {
			cT, _ := parseTime(layOut, time.Now().Format(tstringFormat))
			log.Printf("Stopping Analysis for Today at %+v", cT)

			alive = false
			t.Stop()
			break
		}
		stamp := <-t.C
		log.Printf("Starting Analysis at %+v", stamp.Format(tstringFormat))
		// do actual analysis here
		time.Sleep(time.Second * 3)
		trade.Analyse()
	}
	return nil
}

//Analyse the trade
func (trade *Trade) Analyse() {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", trade.Order.InstrumentToken, trade.Order.TradeInterval)
	orderRespsonse, _ := db.GetTicks()
	var tradeDetailsList []TradeDetails
	for _, results := range orderRespsonse.Results {
		for _, rows := range results.Series {
			td := &TradeDetails{}
			for _, row := range rows.Values {
				td.Close, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[1]), 64)
				td.High, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[2]), 64)
				td.Low, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[3]), 64)
				td.Open, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[4]), 64)
				td.TotalBuyQuantity, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[5]), 64)
				td.TotalSellQuantity, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[6]), 64)
				tradeDetailsList = append(tradeDetailsList, *td)
			}

		}

	}

	log.Printf("Aggregate - %+v", tradeDetailsList)
	log.Println("---------------------------------")

	//return &ordList
}

func getInterval(i string) int {
	switch i {
	case "1m":
		return 1

	case "3m":
		return 3

	case "5m":
		return 5

	default:
		return 0
	}

}

func (trade Trade) stopAnalysis(closingTime time.Time, doWait bool) (bool, error) {
	if doWait {
		time.Sleep(time.Second * 2)
	}
	currentTime := time.Now().Format(tstringFormat)
	parsedCurrentTime, err := parseTime(layOut, currentTime)
	if err != nil {
		return true, err
	}
	return isMarketClose(closingTime, parsedCurrentTime), nil

}

func parseTime(format string, tstring string) (time.Time, error) {
	parsedTime, err := time.Parse(format, tstring)
	if err != nil {
		log.Fatalf("Error parsing market CloseTime Time: %+v", err)
		return time.Time{}, err
	}

	return parsedTime, nil
}

func isMarketClose(closeTime, currentTime time.Time) bool {
	return currentTime.After(closeTime)

}

func waitBeforeAnalysis(interval int) int {
	ct, _ := parseTime(layOut, time.Now().Format(tstringFormat))
	_, min, sec := ct.Clock()
	next := min + (interval - min%interval)
	waitTime := next - min
	if waitTime%interval == 0 {
		return 0
	}
	waitTimeSeconds := (waitTime * 60) - sec
	return waitTimeSeconds

}
