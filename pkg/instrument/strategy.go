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
				SendAlerts(fmt.Sprintf("BUY %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange))
				if cs.Details[1].Low <= cs.getLowestPrice() { //TODO: Use CS data to get lowes price rather than querying DB
					SendAlerts(fmt.Sprintf("BUY %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange))
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
				if cs.Details[1].High > cs.getHighestPrice() { //TODO: Use CS data to get Hishest price rather than querying DB
					SendAlerts(fmt.Sprintf("DEFINITE SELL  %s - %s - %s. %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, "Last day's High Broken!"))
				}
				SendAlerts(fmt.Sprintf("SELL  %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange))
			}

		}

	}
}

// PriceAction strategy
func (cs CandleStick) PriceAction() {
	log.Printf("Instrument: %v", cs.Instrument.Name)
	log.Printf("Previous Trade: %v", cs.PreviousTrade)
	previousDayLow := cs.Details[len(cs.Details)-1].Low
	lowestToday, _ := getLowestLow(cs.Details[:len(cs.Details)-1])
	log.Printf("Previous Day Low: %v", previousDayLow)
	log.Printf("Today's Lowest so far: %v", lowestToday)
	previousDayHigh := cs.Details[len(cs.Details)-1].High
	highestToday, _ := getHighestHigh(cs.Details[:len(cs.Details)-1])
	isBull, bullCount := isBullish(cs.Details)
	isBear, bearCount := isBearish(cs.Details)

	isDozi := isDozi(cs.Details[0])
	bullishMaru := isBullishMarubuzo(cs.Details[0])
	bearishMaru := isBearishMarubuzo(cs.Details[0])
	bullishHammer := isBullishHammer(cs.Details[0])
	bearishHammer := isBearishHammer(cs.Details[0])
	invertedHammer := isInvertedHammer(cs.Details[0])
	shortTrend, shortTrendCount := getShortTermTrend(cs.Details[1:])
	_, bearTrendCount := isBearish(cs.Details[1:])
	_, bullTrendCount := isBullish(cs.Details[1:])
	hhePattern := higherLowsEngulfingPatternCount(cs.Details)
	lhePattern := lowerHighsEngulfingPatternCount(cs.Details)

	if cs.PreviousTrade == "SOLD" || len(cs.PreviousTrade) == 0 {
		if bearishHammer || bullishHammer || isBull || bullishMaru || isDozi {
			if (shortTrend == "decline" && shortTrendCount >= 3) || (bearTrendCount >= 3 || bearCount >= 3) || lhePattern >= 5 {
				if lowestToday > previousDayLow {
					// log.Printf("Instrument: %v", cs.Instrument.Name)
					log.Printf("BUY %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
					log.Printf("Previous Trade: %v :: Bearish Hammer: %v :: bullishHammer: %v :: isBull: %v :: BullishMaru:: %v :: isDozi: %v", cs.PreviousTrade, bearishHammer, bullishHammer, isBull, bullishMaru, isDozi)
					log.Printf("shortTrend: %v :: shortTrendCount: %v :: bearTrendCount: %v :: bearCount: %v :: lhePattern:: %v", shortTrend, shortTrendCount, bearTrendCount, bearCount, lhePattern)
					SendAlerts(fmt.Sprintf("BUY CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange))
				}
			}

		}

	} else if cs.PreviousTrade == "BOUGHT" {
		if isBear || bearishMaru || isDozi {
			if (shortTrend == "rally" && shortTrendCount >= 2) || (bullTrendCount >= 2 || bullCount >= 2) || hhePattern >= 3 {
				// 	log.Printf("Instrument: %v", cs.Instrument.Name)
				log.Printf("SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
				log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
				log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
				SendAlerts(fmt.Sprintf("SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange))
			}

		} else if (shortTrend == "rally" && shortTrendCount >= 3) || (bullTrendCount >= 3 || bullCount >= 3) || hhePattern >= 3 {
			//log.Printf("Instrument: %v", cs.Instrument.Name)
			log.Printf("SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
			log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
			log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
			SendAlerts(fmt.Sprintf("SELL CALL  %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange))
		}
	} else if len(cs.PreviousTrade) == 0 && (highestToday >= previousDayHigh) {
		if isBear || bearishMaru || isDozi || invertedHammer {
			if (shortTrend == "rally" && shortTrendCount >= 2) || (bullTrendCount >= 1 || bullCount >= 1) || hhePattern >= 3 {
				{
					log.Printf("SHORT SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
					log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
					log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
					SendAlerts(fmt.Sprintf("SHORT SELL CALL %s - %s - %s :: MESSAGE: %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, "Ensure that market is falling down"))
				}
			}

		}

	}

	log.Println("-----------------------------------------------------------------------------------------------------------------------")
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
