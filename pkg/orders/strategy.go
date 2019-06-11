package orders

import (
	"log"

	"github.com/stockist/pkg/storage"
)

// BuyLow identifies the lowest price where to stock could be bought.
func (trade Trade) BuyLow() {
	if trade.PreviousTrade == "SOLD" || len(trade.PreviousTrade) == 0 {
		isBull, _ := isBullish(trade.Details)
		dozi := isDozi(trade.Details[0])
		bullishMaru := isBullishMarubuzo(trade.Details[0])
		bullishHammer := isBullishHammer(trade.Details[0])
		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
		bearTrend, bearCounts := isBearish(trade.Details[1:])
		lhePattern := lowerHighsEngulfingPatternCount(trade.Details)
		log.Printf("isBull - %v", isBull)
		log.Printf("Dozi - %v", dozi)
		log.Printf("bullishMaru - %v", bullishMaru)
		log.Printf("bullishHammer - %v", bullishHammer)
		log.Printf("shortTrend - %v", shortTrend)
		log.Printf("trendCount - %v", trendCount)
		log.Printf("bearTrend - %v", bearTrend)
		log.Printf("bearCounts - %v", bearCounts)
		log.Printf("lhePattern - %v", lhePattern)

		if isBull || bullishMaru || bullishHammer || dozi {
			if (shortTrend == "decline" && trendCount >= 3) || (bearTrend && bearCounts >= 3) || lhePattern >= 5 {
				//Good to buy now with stop loss
				// Previous low should be less than open and should be the lowest (and should be lower than last day's low?)
				if trade.Details[1].Low < trade.Details[len(trade.Details)-1].Open && trade.Details[1].Low >= trade.getLowestPrice() {
					log.Print("Best time to buy this stock")
					db := storage.NewDB(DBUrl, StockDB, "trade")
					db.InsertTrade(trade.Order.InstrumentToken, "BUY")
					//Create some alert here
				}
			}

		}

	}
}

// StrategyOne is one of the random strategies to
func (trade Trade) StrategyOne() {

	if trade.PreviousTrade == "SOLD" || len(trade.PreviousTrade) == 0 {
		//Analyse if good to Buy
		isBull, _ := isBullish(trade.Details)
		dozi := isDozi(trade.Details[0])
		bullishMaru := isBullishMarubuzo(trade.Details[0])
		bullishHammer := isBullishHammer(trade.Details[0])
		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
		bearTrend, bearCounts := isBearish(trade.Details[1:])
		lhePattern := lowerHighsEngulfingPatternCount(trade.Details)

		if isBull && (bullishMaru || bullishHammer) {
			if (shortTrend == "decline" && trendCount >= 3) || (bearTrend && bearCounts >= 3) || lhePattern >= 5 {
				//Good to buy now with stop loss
			}
		} else if dozi {
			if (shortTrend == "decline" && trendCount >= 3) || (bearTrend && bearCounts >= 3) || lhePattern >= 5 {
				//Good to buy now with stop loss
			}

		} else if isBull && !bullishMaru {
			if (shortTrend == "decline" && trendCount >= 5) || (bearTrend && bearCounts >= 5) || lhePattern >= 5 {
				//Good to buy now with stop loss
			}
		}

	} else if trade.PreviousTrade == "BOUGHT" {
		//Analyse if good to SELL
		isBear, _ := isBearish(trade.Details)
		//isHammer := isHammer(trade.Details[0])
		dozi := isDozi(trade.Details[0])
		bearishMaru := isBearishMarubuzo(trade.Details[0])
		//shootingStar := isInvertedHammer(trade.Details[0])
		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
		bullTrend, bullCounts := isBullish(trade.Details[1:])
		hlePattern := higherLowsEngulfingPatternCount(trade.Details)

		if isBear || bearishMaru {
			if (shortTrend == "rally" && trendCount >= 3) || (bullTrend && bullCounts >= 3) || hlePattern >= 5 {
				//Good to Sell here
			}
		} else if dozi {
			if (shortTrend == "rally" && trendCount >= 3) || (bullTrend && bullCounts >= 3) || hlePattern >= 5 {
				//Good to Sell here
			}

		} else if isBear && !bearishMaru {
			if (shortTrend == "rally" && trendCount >= 5) || (bullTrend && bullCounts >= 5) || hlePattern >= 5 {
				//Good to Sell here
			}
		}

	}

	log.Println("No action performed on this stock")

}
