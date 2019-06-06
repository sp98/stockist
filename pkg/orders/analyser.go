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

	tradeDetails := trade.getTicks()
	log.Printf("Aggregate results - %+v", tradeDetails)

	/*Start Analyses
	1. is Bullish Candle?
	2. is Bearish Candle?
	3. is Dozi candle?
	4. is marubuzo?
	5. Is hammer?
	6. is inverted hammer?
	7. is shooting star?
	*/
	bullish, bullCount := isBullish(*tradeDetails)
	bearish, bearCount := isBearish(*tradeDetails)
	trend, trendCount := getTrend(*tradeDetails)
	dozi := isDozi((*tradeDetails)[0])
	maru := isMarubuzo((*tradeDetails)[0])
	hammer := isHammer((*tradeDetails)[0])

	log.Println("--- Analysis Results --- ")
	log.Printf("Bullish- %v :: Bullish Count- %v", bullish, bullCount)
	log.Printf("Bearish- %v :: Bearish Count- %v", bearish, bearCount)
	log.Printf("Trend- %v :: Trend Count- %v", trend, trendCount)
	log.Printf("is Dozi? - %v", dozi)
	log.Printf("is Marubuzo- %v", maru)
	log.Printf("is Hammer- %v", hammer)
	log.Println("---------------------------------")

}

func (trade *Trade) getTicks() *[]TradeDetails {
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
	// log.Printf("Aggregate - %+v", tradeDetailsList)
	// log.Println("---------------------------------")

	return &tradeDetailsList
}

func isBullish(tradeDetails []TradeDetails) (bool, int) {

	lastTradeDetail := tradeDetails[0]
	if lastTradeDetail.Open > lastTradeDetail.Close {
		return false, 0
	}

	var trendCount int
	for _, td := range tradeDetails {
		if td.Open > td.Close {
			break
		}
		trendCount = trendCount + 1
	}

	return true, trendCount

}

func isBearish(tradeDetails []TradeDetails) (bool, int) {
	lastTradeDetail := tradeDetails[0]
	if lastTradeDetail.Open < lastTradeDetail.Close {
		return false, 0
	}

	var trendCount int
	for _, td := range tradeDetails {
		if td.Open < td.Close {
			break
		}
		trendCount = trendCount + 1
	}

	return true, trendCount

}

func isDozi(td TradeDetails) bool {
	if td.Open == td.Close && (td.High != td.Open || td.Close != td.Open) {
		return true
	}
	return false

}

func isMarubuzo(td TradeDetails) bool {
	if td.Open == td.Low && td.Close == td.High {
		return true
	}
	return false

}

func isHammer(td TradeDetails) bool {

	if td.Open < td.Close {
		if ((td.Close - td.Open) >= 3*td.High-td.Close) && ((td.Open - td.Low) < (td.Close - td.Open)) {
			return true
		}
	} else if td.Open > td.Close {
		if ((td.Open - td.Close) >= 3*td.High-td.Open) && ((td.Close - td.Low) < (td.Open - td.Close)) {
			return true
		}
	}

	return false
}

func getTrend(tradeDetails []TradeDetails) (string, int) {

	trend := ""
	trendCount := 0

	if len(tradeDetails) < 2 {
		return trend, trendCount
	}

	if tradeDetails[0].High > tradeDetails[1].High {
		trend = "upstrend"
		for i := 0; i < len(tradeDetails)-1; i++ {
			if tradeDetails[i].High < tradeDetails[i+1].High {
				return trend, trendCount
			}
			trendCount = trendCount + 1
		}

	} else if tradeDetails[0].High < tradeDetails[1].High {
		trend = "downtrend"
		for i := 0; i < len(tradeDetails)-1; i++ {
			if tradeDetails[i].High > tradeDetails[i+1].High {
				return trend, trendCount
			}
			trendCount = trendCount + 1
		}
	}

	return "", 0
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
