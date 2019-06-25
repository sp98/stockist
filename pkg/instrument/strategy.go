package instrument

import (
	"fmt"
	"log"

	alerts "github.com/stockist/pkg/notification"
	"github.com/stockist/pkg/storage"
)

var separation = "---------------------------------------------------------"

// PriceAction strategy
func (cs CandleStick) PriceAction() {
	log.Printf("Instrument: %v", cs.Instrument.Name)
	log.Printf("Previous Trade: %v", cs.PreviousTrade)

	prevDayTrend, prevDayChange := getPreviousDayTrend(cs.Details[len(cs.Details)-1].Open, cs.Details[len(cs.Details)-1].Close)
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
		if bearishHammer || bullishHammer || bullishMaru || isDozi || (isBull && !invertedHammer) {
			if (shortTrend == "decline" && shortTrendCount >= 3) || (bearTrendCount >= 3 || bearCount >= 3) || lhePattern >= 5 {
				if lowestToday > previousDayLow {
					log.Printf("BUY %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
					log.Printf("Previous Trade: %v :: Bearish Hammer: %v :: bullishHammer: %v :: isBull: %v :: BullishMaru:: %v :: isDozi: %v", cs.PreviousTrade, bearishHammer, bullishHammer, isBull, bullishMaru, isDozi)
					log.Printf("shortTrend: %v :: shortTrendCount: %v :: bearTrendCount: %v :: bearCount: %v :: lhePattern:: %v", shortTrend, shortTrendCount, bearTrendCount, bearCount, lhePattern)
					msg := fmt.Sprintf("BUY CALL ::  %s - %s - %s\nPrevious Day Trend: %s \nPrevious Day Change: %s \nMESSAGE : %s \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, "Ensure that market is moving up", separation)
					alerts.SendAlerts(msg, alerts.BuyStockChannel)
				}
			}

		}

		if len(cs.PreviousTrade) == 0 {
			if invertedHammer || bearishMaru || isDozi || isBear {
				if (shortTrend == "rally" && shortTrendCount >= 3) || (bullTrendCount >= 3 || bullCount >= 3) || hhePattern >= 5 {
					if highestToday <= previousDayHigh && lowestToday <= previousDayLow {
						log.Printf("SHORT SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
						log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
						log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
						msg := fmt.Sprintf("SHORT SELL CALL :: %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s \nMESSAGE: %s  \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, "Ensure that market is falling down", separation)
						alerts.SendAlerts(msg, alerts.ShortSellStocksChannel)

					}
				}

			}

		}

	} else if cs.PreviousTrade == "BOUGHT" {
		if isBear || bearishMaru || isDozi {
			if (shortTrend == "rally" && shortTrendCount >= 2) || (bullTrendCount >= 2 || bullCount >= 2) {
				log.Printf("SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
				log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
				log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
				msg := fmt.Sprintf("SELL CALL %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, separation)
				alerts.SendAlerts(msg, alerts.SellStockChannel)
			}

		}
	}

	log.Println("-----------------------------------------------------------------------------------------------------------------------")
}

//OpeningTrend checks whether the opened low or open high and gives the change percent as well.
func (cs CandleStick) OpeningTrend() {
	name := cs.Instrument.Name
	symbol := cs.Instrument.Symbol
	exchange := cs.Instrument.Exchange
	//log.Printf("Instrument: %v", name)
	var openPrice float64
	for i := len(cs.Details) - 1; i >= 0; i-- {
		// log.Println(cs.Details[i].AverageTradedPrice)
		if cs.Details[i].AverageTradedPrice != 0 {
			openPrice = cs.Details[i].AverageTradedPrice //TODO: for some reasons open price is the average trade price!
			break
		}
	}

	if openPrice == 0 {
		log.Printf("Couldn't find open price for the instruments - %v", cs.Instrument.Name)
		return
	}

	//prevDayLow := cs.Details[len(cs.Details)-1].Low
	//prevDayHigh := cs.Details[len(cs.Details)-1].High
	prevDayClose := cs.Details[len(cs.Details)-1].Close
	//prevDayOpen := cs.Details[len(cs.Details)-1].Open

	// if openPrice < prevDayLow {
	// 	//Today's open is lower than previous day's low
	// 	change := ((prevDayLow - openPrice) / prevDayLow) * 100
	// 	SendAlerts(fmt.Sprintf("OPENING TRADE :: %s - %s - %s \nOpen: %.2f \nPrevious Low: %.2f \nChange: -%.2f%% \n%s", name, symbol, exchange, openPrice, prevDayLow, change, separation))
	// }else if openPrice >= prevDayHigh {
	// 	change := ((openPrice - prevDayHigh) / openPrice) * 100
	// 	//Today's open is lower than previous day's High
	// 	SendAlerts(fmt.Sprintf("OPENING TRADE :: %s - %s - %s \nOpen: %.2f \nPrevious High: %.2f \nChange: +%.2f%% \n%s", name, symbol, exchange, openPrice, prevDayHigh, change, separation))
	// }

	if openPrice < prevDayClose {
		//Today's open is lower than previous day's Open
		change := ((prevDayClose - openPrice) / prevDayClose) * 100
		msg := fmt.Sprintf("OPENING TRADE :: %s - %s - %s \nOpen: %.2f \nPrevious Close: %.2f \nChange: -%.2f%% \n%s", name, symbol, exchange, openPrice, prevDayClose, change, separation)
		alerts.SendAlerts(msg, alerts.OpenTrendChannel)
	} else if openPrice > prevDayClose {
		//Today's open is greater than previous day's open
		change := ((openPrice - prevDayClose) / openPrice) * 100
		msg := fmt.Sprintf("OPENING TRADE :: %s - %s - %s \nOpen: %.2f \nPrevious Close: %.2f \nChange: +%.2f%% \n%s", name, symbol, exchange, openPrice, prevDayClose, change, separation)
		alerts.SendAlerts(msg, alerts.OpenTrendChannel)
	}

}

//AnalyseSensex analyses the rally and declines in Sensex
func (cs CandleStick) AnalyseSensex() {
	log.Printf("Instrument: %v", cs.Instrument.Name)
	previousDayLow := cs.Details[len(cs.Details)-2].Low
	lowestToday, _ := getLowestLow(cs.Details[:len(cs.Details)-1])
	previousDayHigh := cs.Details[len(cs.Details)-2].High
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
	hhePatternTrend := higherLowsEngulfingPatternCount(cs.Details[1:])
	lhePattern := lowerHighsEngulfingPatternCount(cs.Details)          ////Lower High Pattern including last CS
	lhePatternTrend := lowerHighsEngulfingPatternCount(cs.Details[1:]) //Lower High Pattern count before last CS

	if bearishHammer || bullishHammer || bullishMaru || isDozi || (isBull && !invertedHammer) {
		if (shortTrend == "decline" && shortTrendCount >= 3) || (bearTrendCount >= 3 || bearCount >= 3) || lhePattern >= 4 || lhePatternTrend >= 4 {
			if lowestToday > previousDayLow {
				msg := fmt.Sprintf("%s \n%s", "ALERT: Sensex declined below previous day's low. Look to buy now", separation)
				alerts.SendAlerts(msg, alerts.SensexTrendChannel)
			} else {
				msg := fmt.Sprintf("%s \n%s", "ALERT: Sensex just declined. Look to buy now", separation)
				alerts.SendAlerts(msg, alerts.SensexTrendChannel)
			}

		}
	} else if invertedHammer || bearishMaru || isDozi || isBear {
		if (shortTrend == "rally" && shortTrendCount >= 3) || (bullTrendCount >= 3 || bullCount >= 3) || hhePattern >= 3 || hhePatternTrend >= 4 {
			if highestToday >= previousDayHigh {
				msg := fmt.Sprintf("%s \n%s", "ALERT: Sensex just rallied above previous day's high. Look to sell now", separation)
				alerts.SendAlerts(msg, alerts.SensexTrendChannel)
			} else {
				msg := fmt.Sprintf("%s \n%s", "ALERT: Sensex just rallied. Look to sell now", separation)
				alerts.SendAlerts(msg, alerts.SensexTrendChannel)
			}

		}

	}

}

func getPreviousDayTrend(open, close float64) (string, string) {
	var trend string
	var change string
	if open < close {
		trend = "uptrend"
		change = fmt.Sprintf("+%.2f%%", ((close-open)/close)*100)
	} else if open > close {
		trend = "downtrend"
		change = fmt.Sprintf("-%.2f%%", ((open-close)/open)*100)
	} else {
		trend = "sideways"
		change = "0.00"
	}

	return trend, change

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

//Logs

// log.Printf("isBull - %v", isBull)
// log.Printf("isBear - %v", isBear)
// log.Printf("bearCount - %v", bearCount)
// log.Printf("Dozi - %v", isDozi)
// log.Printf("bullishMaru - %v", bullishMaru)
// log.Printf("bearishMaru - %v", bearishMaru)
// log.Printf("bullishHammer - %v", bullishHammer)
// log.Printf("bearishHammer - %v", bearishHammer)
// log.Printf("shortTrend - %v", shortTrend)
// log.Printf("bearTrendCount - %v", bearTrendCount)
// log.Printf("bullTrendCount - %v", bullTrendCount)
// log.Printf("isInvertedHammer - %v", invertedHammer)
// log.Printf("hhePattern - %v", hhePattern)
// log.Printf("hhePatternTrend - %v", hhePatternTrend)
// log.Printf("lhePatternTrend - %v", lhePatternTrend)
