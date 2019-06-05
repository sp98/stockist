package orders

import (
	"fmt"
	"log"
	"time"
)

var (
	marketCloseTime = "%sT14:38:00.000Z"
	tstringFormat   = "2006-01-02T15:04:07.000Z"
)

//Trade holds the currently trading order details
type Trade struct {
	Order             Order
	IsBullish         bool
	CurrentTrend      string
	CurrentTrendCount int
}

func (trade *Trade) startAnalysis() error {
	log.Println("--- Begin Analysis ----")

	mct, err := parseTime(time.RFC3339, fmt.Sprintf(marketCloseTime, trade.Order.TradeDate))
	if err != nil {
		return err
	}

	if stop, _ := trade.stopAnalysis(mct); stop {
		log.Printf("Can't start analysis. Already past market closing time %+v", mct)
		return nil
	}

	wt := waitBeforeAnalysis()
	if wt > 0 {
		log.Printf("Wait for %d minutes before starting", wt)
		time.Sleep(time.Second * time.Duration(wt))
	}

	interval := getInterval(trade.Order.TradeInterval)
	if interval == 0 {
		log.Fatal("invalid order interval ")
		//return error here
	}

	t := time.NewTicker(time.Minute * time.Duration(interval))

	log.Printf("Analysis Start Time: %+v ::: Analysis Stop Time: %+v", time.Now(), fmt.Sprintf(marketCloseTime, trade.Order.TradeDate))

	for alive := true; alive; {
		res, _ := trade.stopAnalysis(mct)
		if res {
			cT, _ := parseTime(time.RFC3339, time.Now().Format(tstringFormat))
			log.Printf("Stopping Analysis for Today at %+v", cT)

			alive = false
			t.Stop()
			break
		}
		stamp := <-t.C
		log.Printf("Starting Analysis at %+v", stamp.Format(tstringFormat))
		// do actual analysis here
	}
	return nil
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

func (trade Trade) stopAnalysis(closingTime time.Time) (bool, error) {
	currentTime := time.Now().Format(tstringFormat)
	parsedCurrentTime, err := parseTime(time.RFC3339, currentTime)
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

func waitBeforeAnalysis() int {
	ct, _ := parseTime(time.RFC3339, time.Now().Format(tstringFormat))
	_, min, sec := ct.Clock()
	log.Printf("Mins - %d", min)
	log.Printf("sec - %d", sec)
	next := min + (5 - min%5)
	waitTime := next - min
	if waitTime%5 == 0 {
		return 0
	}
	waitTimeSeconds := (waitTime * 60) + sec
	return waitTimeSeconds

}
