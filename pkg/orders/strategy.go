package orders

import (
	"fmt"
	"log"

	"github.com/stockist/pkg/storage"
)

// BuyLow identifies the lowest price where to stock could be bought.
func (trade Trade) BuyLow() {

	lastTrade := getLastTrade(trade.Order.InstrumentToken)
	log.Printf("Last Trade : %v", lastTrade)
	if lastTrade == "SOLD" || len(lastTrade) == 0 {
		isBull, bullTrend := isBullish(trade.Details)
		dozi := isDozi(trade.Details[0])
		bullishMaru := isBullishMarubuzo(trade.Details[0])
		bearishMaru := isBearishMarubuzo(trade.Details[0])
		bullishHammer := isBullishHammer(trade.Details[0])
		bearishHammer := isBearishHammer(trade.Details[0])
		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
		bearTrend, bearCounts := isBearish(trade.Details[1:])
		lhePattern := lowerHighsEngulfingPatternCount(trade.Details)
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
				if trade.Details[1].Low <= trade.getLowestPrice() {
					//Create some alert here
					SendAlerts(fmt.Sprintf("BUY %s-%s", trade.Order.InstrumentName, trade.Order.Exchange))
				}
			}

		}

	} else if lastTrade == "BOUGHT" {
		isBear, bearCount := isBearish(trade.Details)
		dozi := isDozi(trade.Details[0])
		bullishMaru := isBullishMarubuzo(trade.Details[0])
		bearishMaru := isBearishMarubuzo(trade.Details[0])
		bullishHammer := isBullishHammer(trade.Details[0])
		bearishHammer := isBearishHammer(trade.Details[0])
		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
		bearTrend, bearCounts := isBearish(trade.Details[1:])
		bullTrend, bullCounts := isBullish(trade.Details[1:])
		hhePattern := higherLowsEngulfingPatternCount(trade.Details)
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
				//Good to SELL now with stop loss
				// Last High should be higher than last days high.
				// Last High should
				SendAlerts(fmt.Sprintf("SELL %s", trade.Order.InstrumentName))

				if trade.Details[1].High > trade.getHighestPrice() {
					//updateTradeInDB(trade.Order.InstrumentToken, "BUY")
					//Create some alert here
					SendAlerts(fmt.Sprintf("DEFINITE SELL %s", trade.Order.InstrumentName))
				}
			}

		}

	}
}

func getLastTrade(instToken string) string {
	db := storage.NewDB(DBUrl, StockDB, "trade")
	trade, _ := db.GetLastTrade(instToken)
	return trade

}

// // StrategyOne is one of the random strategies to
// func (trade Trade) StrategyOne() {

// 	if trade.PreviousTrade == "SOLD" || len(trade.PreviousTrade) == 0 {
// 		//Analyse if good to Buy
// 		isBull, _ := isBullish(trade.Details)
// 		dozi := isDozi(trade.Details[0])
// 		bullishMaru := isBullishMarubuzo(trade.Details[0])
// 		bullishHammer := isBullishHammer(trade.Details[0])
// 		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
// 		bearTrend, bearCounts := isBearish(trade.Details[1:])
// 		lhePattern := lowerHighsEngulfingPatternCount(trade.Details)

// 		if isBull && (bullishMaru || bullishHammer) {
// 			if (shortTrend == "decline" && trendCount >= 3) || (bearTrend && bearCounts >= 3) || lhePattern >= 5 {
// 				//Good to buy now with stop loss
// 			}
// 		} else if dozi {
// 			if (shortTrend == "decline" && trendCount >= 3) || (bearTrend && bearCounts >= 3) || lhePattern >= 5 {
// 				//Good to buy now with stop loss
// 			}

// 		} else if isBull && !bullishMaru {
// 			if (shortTrend == "decline" && trendCount >= 5) || (bearTrend && bearCounts >= 5) || lhePattern >= 5 {
// 				//Good to buy now with stop loss
// 			}
// 		}

// 	} else if trade.PreviousTrade == "BOUGHT" {
// 		//Analyse if good to SELL
// 		isBear, _ := isBearish(trade.Details)
// 		//isHammer := isHammer(trade.Details[0])
// 		dozi := isDozi(trade.Details[0])
// 		bearishMaru := isBearishMarubuzo(trade.Details[0])
// 		//shootingStar := isInvertedHammer(trade.Details[0])
// 		shortTrend, trendCount := getShortTermTrend(trade.Details[1:])
// 		bullTrend, bullCounts := isBullish(trade.Details[1:])
// 		hlePattern := higherLowsEngulfingPatternCount(trade.Details)

// 		if isBear || bearishMaru {
// 			if (shortTrend == "rally" && trendCount >= 3) || (bullTrend && bullCounts >= 3) || hlePattern >= 5 {
// 				//Good to Sell here
// 			}
// 		} else if dozi {
// 			if (shortTrend == "rally" && trendCount >= 3) || (bullTrend && bullCounts >= 3) || hlePattern >= 5 {
// 				//Good to Sell here
// 			}

// 		} else if isBear && !bearishMaru {
// 			if (shortTrend == "rally" && trendCount >= 5) || (bullTrend && bullCounts >= 5) || hlePattern >= 5 {
// 				//Good to Sell here
// 			}
// 		}

// 	}

// 	log.Println("No action performed on this stock")

// }

// func sendMail(body string) {
// 	from := "spillai1098@gmail.com"
// 	pass := "$ankump1l"
// 	to := "sapillai@redhat.com"

// 	msg := "From: " + from + "\n" +
// 		"To: " + to + "\n" +
// 		"Subject: IMPORTANT: STOCKIST\n\n" +
// 		body

// 	err := smtp.SendMail("smtp.gmail.com:587",
// 		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
// 		from, []string{to}, []byte(msg))

// 	if err != nil {
// 		log.Printf("smtp error: %s", err)
// 		return
// 	}

// 	log.Print("sent")
// }
