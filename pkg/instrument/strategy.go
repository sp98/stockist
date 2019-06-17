package instrument

import (
	"fmt"
	"log"

	"github.com/stockist/pkg/storage"
)

// BuyLowSellHigh helps to buy the stock at lowest price and sell at highest.
func (cs CandleStick) BuyLowSellHigh() {

	lastTrade := getLastTrade(cs.Instrument.Token)
	log.Printf("Last Trade : %v", lastTrade)
	if lastTrade == "SOLD" || len(lastTrade) == 0 {
		isBull, bullTrend := isBullish(cs.Details)
		dozi := isDozi(cs.Details[0])
		bullishMaru := isBullishMarubuzo(cs.Details[0])
		bearishMaru := isBearishMarubuzo(cs.Details[0])
		bullishHammer := isBullishHammer(cs.Details[0])
		bearishHammer := isBearishHammer(cs.Details[0])
		shortTrend, trendCount := getShortTermTrend(cs.Details[1:])
		bearTrend, bearCounts := isBearish(cs.Details[1:])
		lhePattern := lowerHighsEngulfingPatternCount(cs.Details)
		log.Printf("isBull - %v", isBull)
		log.Printf("bullTrend - %v", bullTrend)
		log.Printf("Dozi - %v", dozi)
		log.Printf("bullishMaru - %v", bullishMaru)
		log.Printf("bearishMaru - %v", bearishMaru)
		log.Printf("bullishHammer - %v", bullishHammer)
		log.Printf("bearishHammer - %v", bearishHammer)
		log.Printf("shortTrend - %v", shortTrend)
		log.Printf("trendCount - %v", trendCount)
		log.Printf("bearTrend - %v", bearTrend)
		log.Printf("bearCounts - %v", bearCounts)
		log.Printf("lhePattern - %v", lhePattern)

		if isBull || bullishMaru || bullishHammer || dozi {
			if (shortTrend == "decline" && trendCount >= 3) || (bearTrend && bearCounts >= 3) || lhePattern >= 5 {
				//Good to buy now with stop loss
				// Previous low should be the lowest so far. Lower than both today's low and previous day's low.
				if cs.Details[1].Low <= cs.getLowestPrice() {
					//Create some alert here
					SendAlerts(fmt.Sprintf("BUY %s-%s-%s", cs.Instrument.Name, cs.Instrument.Token, cs.Instrument.Exchange))
				}
			}

		}

	} else if lastTrade == "BOUGHT" {
		isBear, bearCount := isBearish(cs.Details)
		dozi := isDozi(cs.Details[0])
		bullishMaru := isBullishMarubuzo(cs.Details[0])
		bearishMaru := isBearishMarubuzo(cs.Details[0])
		bullishHammer := isBullishHammer(cs.Details[0])
		bearishHammer := isBearishHammer(cs.Details[0])
		shortTrend, trendCount := getShortTermTrend(cs.Details[1:])
		bearTrend, bearCounts := isBearish(cs.Details[1:])
		bullTrend, bullCounts := isBullish(cs.Details[1:])
		hhePattern := higherLowsEngulfingPatternCount(cs.Details)
		log.Printf("isBear - %v", isBear)
		log.Printf("bearCount - %v", bearCount)
		log.Printf("Dozi - %v", dozi)
		log.Printf("bullishMaru - %v", bullishMaru)
		log.Printf("bearishMaru - %v", bearishMaru)
		log.Printf("bullishHammer - %v", bullishHammer)
		log.Printf("bearishHammer - %v", bearishHammer)
		log.Printf("shortTrend - %v", shortTrend)
		log.Printf("trendCount - %v", trendCount)
		log.Printf("bearTrend - %v", bearTrend)
		log.Printf("bearCounts - %v", bearCounts)
		log.Printf("hhePattern - %v", hhePattern)

		if isBear || bearishMaru || dozi {
			if (shortTrend == "rally" && trendCount >= 3) || (bullTrend && bullCounts >= 3) || hhePattern >= 5 {
				if cs.Details[1].High > cs.getHighestPrice() {
					SendAlerts(fmt.Sprintf("SELL  %s-%s-%s", cs.Instrument.Name, cs.Instrument.Token, cs.Instrument.Exchange))
				}
			}

		}

	}
}

func updateTradeInDB(option, instToken string) {
	db := storage.NewDB(DBUrl, StockDB, "trade")
	db.InsertTrade(instToken, option)

}

func getLastTrade(instToken string) string {
	db := storage.NewDB(DBUrl, StockDB, "trade")
	trade, _ := db.GetLastTrade(instToken)
	return trade

}
