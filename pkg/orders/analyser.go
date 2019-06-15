package orders

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stockist/pkg/storage"
)

var (
	marketCloseTime      = "%s 15:30:00"
	marketActualOpenTime = "%s 09:13:00 MST"
	tstringFormat        = "2006-01-02 15:04:05"
	layOut               = "2006-01-02 15:04:05"
	influxLayout         = "2006-01-02T15:04:05Z"
)

//Trade holds the currently trading order details
type Trade struct {
	Order         Order
	Details       []TradeDetails
	PreviousTrade string
}

//TradeDetails is the aggregate of the trade
type TradeDetails struct {
	AverageTradedPrice float64
	Close              float64
	High               float64
	Low                float64
	Open               float64
	TotalBuyQuantity   float64
	TotalSellQuantity  float64
}

func (trade *Trade) startAnalysis() error {
	log.Printf("--- Begin Analysis for %s ----", trade.Order.InstrumentName)

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

	//trade.Analyse()
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

	if len(*tradeDetails) == 0 {
		log.Println("Error: Trade details is empty!")
		return
	}
	trade.Details = *tradeDetails

	//trade.StrategyOne()
	if len(trade.Details) > 3 {
		trade.BuyLow()
	}
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
				td.AverageTradedPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[1]), 64)
				td.Close, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[2]), 64)
				td.High, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[3]), 64)
				td.Low, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[4]), 64)
				td.Open, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[5]), 64)
				td.TotalBuyQuantity, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[6]), 64)
				td.TotalSellQuantity, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[7]), 64)
				tradeDetailsList = append(tradeDetailsList, *td)
			}

		}

	}

	return &tradeDetailsList
}

func (trade *Trade) getOverallTrend(currentHigh float64) string {

	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", trade.Order.InstrumentToken, trade.Order.TradeInterval)
	maxHigh, _ := db.GetMaxHigh()
	if maxHigh > currentHigh {
		return "downtrend" // downtrend: if current high is not the maximum high.
	}

	return "uptrend" // Uptrend: if current High is the maximum High

}

//Gives the trend before the current Candlestick pattern
func getShortTermTrend(tradeDetails []TradeDetails) (string, int) {
	trend := ""
	trendCount := 0

	if len(tradeDetails) < 2 {
		return trend, trendCount
	}

	if tradeDetails[0].High > tradeDetails[1].High && tradeDetails[0].Low > tradeDetails[1].Low {
		trend = "rally"
		for i := 0; i < len(tradeDetails)-1; i++ {
			if tradeDetails[i].High > tradeDetails[i+1].High && tradeDetails[i].Low > tradeDetails[i+1].Low {
				trendCount = trendCount + 1
				continue
			}
			return trend, trendCount
		}
		return trend, trendCount

	} else if tradeDetails[0].High < tradeDetails[1].High && tradeDetails[0].Low < tradeDetails[1].Low {
		trend = "decline"
		for i := 0; i < len(tradeDetails)-1; i++ {
			if tradeDetails[i].High < tradeDetails[i+1].High && tradeDetails[i].Low < tradeDetails[i+1].Low {
				trendCount = trendCount + 1
				continue
			}
			return trend, trendCount
		}
		return trend, trendCount

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

func getActualMarketOpenTime(date string) (string, error) {
	IST, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println(err)
	}

	mast := fmt.Sprintf(marketActualOpenTime, date)
	longForm := "2006-01-02 15:04:05 MST"
	t, err := time.ParseInLocation(longForm, mast, IST)
	if err != nil {
		return "", err
	}
	return t.UTC().Format("2006-01-02T15:04:05Z"), nil

}
