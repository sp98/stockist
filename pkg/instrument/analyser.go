package instrument

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stockist/pkg/storage"
)

//CandleStick holds the currently trading order details
type CandleStick struct {
	Instrument    Instrument
	Details       []CandleStickList
	PreviousTrade string
}

//CandleStickList is the aggregate of the trade
type CandleStickList struct {
	AverageTradedPrice float64
	Close              float64
	High               float64
	Low                float64
	Open               float64
}

func (cs *CandleStick) startAnalysis() error {
	log.Printf("---ANALYSE - %s ----", cs.Instrument.Name)
	interval := getInterval(cs.Instrument.Interval)
	if interval == 0 {
		log.Fatal("invalid order interval ")
		//return error here
	}

	// mct, err := parseTime(layOut, fmt.Sprintf(marketCloseTime, getDate()))
	// if err != nil {
	// 	return err
	// }

	// if stop, _ := cs.stopAnalysis(mct, false); stop {
	// 	log.Printf("Can't start analysis. Already past market closing time %+v", mct)
	// 	return nil
	// }

	wt := waitBeforeAnalysis(interval)
	if wt > 0 {
		log.Printf("Wait for %.2f minutes before starting", float64(wt)/60)
		time.Sleep(time.Second * time.Duration(wt))
	}

	//cs.Analyse()
	t := time.NewTicker(time.Minute * time.Duration(interval))

	log.Printf("Analysis Start Time: %+v ::: Analysis Stop Time: %+v", time.Now(), fmt.Sprintf(marketCloseTime, getDate()))

	for alive := true; alive; {
		// res, _ := cs.stopAnalysis(mct, true)
		// if res {
		// 	cT, _ := parseTime(layOut, time.Now().Format(tstringFormat))
		// 	log.Printf("Stopping Analysis for Today at %+v", cT)

		// 	alive = false
		// 	t.Stop()
		// 	break
		// }
		stamp := <-t.C
		log.Printf("Starting Analysis at %+v", stamp.Format(tstringFormat))
		// do actual analysis here
		time.Sleep(time.Second * 3)
		cs.Analyse()
	}
	return nil
}

//Analyse the trade
func (cs *CandleStick) Analyse() {

	csList := cs.getTicks()
	log.Printf("Aggregate results - %+v", csList)

	if len(*csList) == 0 {
		log.Println("Error: Candle Stick details are empty!")
		return
	}
	cs.Details = *csList

	if len(cs.Details) > 3 {
		cs.BuyLowSellHigh()
	}
}

func (cs *CandleStick) getTicks() *[]CandleStickList {
	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", cs.Instrument.Token, cs.Instrument.Interval)
	orderRespsonse, _ := db.GetTicks()
	var csDetailList []CandleStickList
	for _, results := range orderRespsonse.Results {
		for _, rows := range results.Series {
			cs := &CandleStickList{}
			for _, row := range rows.Values {
				cs.AverageTradedPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[1]), 64)
				cs.Close, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[2]), 64)
				cs.High, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[3]), 64)
				cs.Low, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[4]), 64)
				cs.Open, _ = strconv.ParseFloat(fmt.Sprintf("%v", row[5]), 64)
				csDetailList = append(csDetailList, *cs)
			}
		}

	}

	return &csDetailList
}

func (cs *CandleStick) getOverallTrend(currentHigh float64) string {

	db := storage.NewDB(DBUrl, StockDB, "")
	db.Measurement = fmt.Sprintf("%s_%s_%s", "ticks", cs.Instrument.Token, cs.Instrument.Interval)
	maxHigh, _ := db.GetMaxHigh()
	if maxHigh > currentHigh {
		return "downtrend" // downtrend: if current high is not the maximum high.
	}

	return "uptrend" // Uptrend: if current High is the maximum High

}

//Gives the trend before the current Candlestick pattern
func getShortTermTrend(csDetails []CandleStickList) (string, int) {
	trend := ""
	trendCount := 0

	if len(csDetails) < 2 {
		return trend, trendCount
	}

	if csDetails[0].High > csDetails[1].High && csDetails[0].Low > csDetails[1].Low {
		trend = "rally"
		for i := 0; i < len(csDetails)-1; i++ {
			if csDetails[i].High > csDetails[i+1].High && csDetails[i].Low > csDetails[i+1].Low {
				trendCount = trendCount + 1
				continue
			}
			return trend, trendCount
		}
		return trend, trendCount

	} else if csDetails[0].High < csDetails[1].High && csDetails[0].Low < csDetails[1].Low {
		trend = "decline"
		for i := 0; i < len(csDetails)-1; i++ {
			if csDetails[i].High < csDetails[i+1].High && csDetails[i].Low < csDetails[i+1].Low {
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

func (cs CandleStick) stopAnalysis(closingTime time.Time, doWait bool) (bool, error) {
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

func getDate() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}