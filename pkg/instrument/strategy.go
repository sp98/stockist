package instrument

import (
	"fmt"
	"log"
	"time"

	kiteticker "github.com/sp98/gokiteconnect/ticker"
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

//OpenLowHigh strategy for stock recommendation
func (cs CandleStick) OpenLowHigh() (string, error) {

	// if len(cs.Details) < 6 {
	// 	log.Println("No Open Low High recommendations before 9:45 am")
	// 	return "", nil
	// }

	if cs.KC == nil {
		log.Println("Error: Kite connection is nil")
		return "", nil
	}

	open, high, low, err := cs.GetOHLC()
	if err != nil {
		log.Printf("Error Finding OHLC for %s. Error : %+v", cs.Instrument.Symbol, err)
		return "", err
	}

	bq, sq, qchange, _ := cs.GetTradeQuantity()

	//Send alert about Open=High and Open=Low stocks. Unsubscribe and stop analysis of the stocks that don't follow Open High Low
	if len(cs.Details) == 4 || len(cs.Details) == 6 || len(cs.Details) == 9 || len(cs.Details) == 12 { //Open == high is a good canditate to short cell in case of negative markets.
		if open == high {
			msg := fmt.Sprintf("Possible Short Sell Stock in downtrend \nInstrument: %s \n Open: %.2f \nHigh: %.2f \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v \n%s", cs.Instrument.Symbol, open, high, bq, sq, qchange, separation)
			alerts.SendAlerts(msg, alerts.OpenLowHigh)
		} else if open == low { //Open == low is a good canditate to buy in case of positive markets.
			msg := fmt.Sprintf("Possible Buy Stock in Uptrend \nInstrument: %s \n Open: %.2f \nLow: %.2f \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v, \n%s", cs.Instrument.Symbol, open, low, bq, sq, qchange, separation)
			alerts.SendAlerts(msg, alerts.OpenLowHigh)
		} else {
			// err := cs.UnSubcribe()
			// if err != nil {
			// 	msg := fmt.Sprintf("Error unsubscribing stock %s", cs.Instrument.Symbol)
			// 	log.Println(msg)
			// 	alerts.SendAlerts(msg, alerts.ErrorChannel)
			// }
			return StopAnalysis, nil
		}
	}

	log.Printf("Instrument: %v", cs.Instrument.Name)
	log.Printf("Previous Trade: %v", cs.PreviousTrade)

	prevDayTrend, prevDayChange := getPreviousDayTrend(cs.Details[len(cs.Details)-1].Open, cs.Details[len(cs.Details)-1].Close)
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

		if open == low { // Buy Call in Uptrend Market
			if bearishHammer || bullishHammer || bullishMaru || isDozi || (isBull && !invertedHammer) {
				if len(cs.Details) <= 18 { //Giving suggestions before 10:30 am
					if (shortTrend == "decline" && shortTrendCount >= 2) || (bearTrendCount >= 2 || bearCount >= 2) || lhePattern >= 3 {
						log.Printf("BUY %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
						log.Printf("Previous Trade: %v :: Bearish Hammer: %v :: bullishHammer: %v :: isBull: %v :: BullishMaru:: %v :: isDozi: %v", cs.PreviousTrade, bearishHammer, bullishHammer, isBull, bullishMaru, isDozi)
						log.Printf("shortTrend: %v :: shortTrendCount: %v :: bearTrendCount: %v :: bearCount: %v :: lhePattern:: %v", shortTrend, shortTrendCount, bearTrendCount, bearCount, lhePattern)
						msg := fmt.Sprintf("BUY CALL ::  %s - %s - %s\nPrevious Day Trend: %s \nPrevious Day Change: %s \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v  \nMESSAGE : %s \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, "Ensure that market is moving up", separation)
						alerts.SendAlerts(msg, alerts.BuyStockChannel)
					}
				} else {
					if (shortTrend == "decline" && shortTrendCount >= 4) || (bearTrendCount >= 4 || bearCount >= 4) || lhePattern >= 5 {
						log.Printf("BUY %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
						log.Printf("Previous Trade: %v :: Bearish Hammer: %v :: bullishHammer: %v :: isBull: %v :: BullishMaru:: %v :: isDozi: %v", cs.PreviousTrade, bearishHammer, bullishHammer, isBull, bullishMaru, isDozi)
						log.Printf("shortTrend: %v :: shortTrendCount: %v :: bearTrendCount: %v :: bearCount: %v :: lhePattern:: %v", shortTrend, shortTrendCount, bearTrendCount, bearCount, lhePattern)
						msg := fmt.Sprintf("BUY CALL ::  %s - %s - %s\nPrevious Day Trend: %s \nPrevious Day Change: %s \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v  \nMESSAGE : %s \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, "Ensure that market is moving up", separation)
						alerts.SendAlerts(msg, alerts.BuyStockChannel)
					}
				}

			}
		}

		if open == high { // SHORT SELL Call in Downtrend Market
			if len(cs.PreviousTrade) == 0 {
				if invertedHammer || bearishMaru || isDozi || isBear {
					if len(cs.Details) <= 18 {
						if (shortTrend == "rally" && shortTrendCount >= 2) || (bullTrendCount >= 2 || bullCount >= 2) || hhePattern >= 4 {
							log.Printf("SHORT SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
							log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
							log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
							msg := fmt.Sprintf("SHORT SELL CALL :: %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v \nMESSAGE: %s  \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, "Ensure that market is falling down", separation)
							alerts.SendAlerts(msg, alerts.ShortSellStocksChannel)

						} else {
							if (shortTrend == "rally" && shortTrendCount >= 4) || (bullTrendCount >= 4 || bullCount >= 4) || hhePattern >= 4 {
								log.Printf("SHORT SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
								log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
								log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
								msg := fmt.Sprintf("SHORT SELL CALL :: %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v  \nMESSAGE: %s  \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, "Ensure that market is falling down", separation)
								alerts.SendAlerts(msg, alerts.ShortSellStocksChannel)

							}
						}
					}
				}

			}
		}

	} else if cs.PreviousTrade == "BOUGHT" {
		if isBear || bearishMaru || isDozi {
			if (shortTrend == "rally" && shortTrendCount >= 4) || (bullTrendCount >= 4 || bullCount >= 4) {
				log.Printf("SELL CALL %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
				log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
				log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
				msg := fmt.Sprintf("SELL CALL %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s  \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v  \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, separation)
				alerts.SendAlerts(msg, alerts.SellStockChannel)
			}

		}
	}

	log.Println("-----------------------------------------------------------------------------------------------------------------------")

	return "", nil

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

	prevDayClose := cs.Details[len(cs.Details)-1].Close

	if openPrice < prevDayClose {
		//Today's open is lower than previous day's Open
		change := ((prevDayClose - openPrice) / prevDayClose) * 100
		// if change >= 1.5 {
		msg := fmt.Sprintf("OPENING TRADE :: %s - %s - %s \nOpen: %.2f \nPrevious Close: %.2f \nChange: -%.2f%% \n%s", name, symbol, exchange, openPrice, prevDayClose, change, separation)
		alerts.SendAlerts(msg, alerts.OpenTrendChannel)
		// }

	} else if openPrice > prevDayClose {
		//Today's open is greater than previous day's open
		change := ((openPrice - prevDayClose) / openPrice) * 100
		// if change >= 1.5 {
		msg := fmt.Sprintf("OPENING TRADE :: %s - %s - %s \nOpen: %.2f \nPrevious Close: %.2f \nChange: +%.2f%% \n%s", name, symbol, exchange, openPrice, prevDayClose, change, separation)
		alerts.SendAlerts(msg, alerts.OpenTrendChannel)
		// }
	}

}

//AnalyseBajajFinance analyses Bajaj Finance Stock
func (cs CandleStick) AnalyseBajajFinance() (string, error) {

	if cs.KC == nil {
		log.Println("Error: Kite connection is nil")
		return "", nil
	}

	log.Printf("Instrument: %v", cs.Instrument.Name)
	log.Printf("Previous Trade: %v", cs.PreviousTrade)

	prevDayTrend, prevDayChange := getPreviousDayTrend(cs.Details[len(cs.Details)-1].Open, cs.Details[len(cs.Details)-1].Close)
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

	bq, sq, qchange, _ := cs.GetTradeQuantity()

	//Send alert about Open=High and Open=Low stocks. Unsubscribe and stop analysis of the stocks that don't follow Open High Low
	if len(cs.Details) > 4 && len(cs.Details) < 15 { //Open == high is a good canditate to short cell in case of negative markets.
		if cs.PreviousTrade == "SOLD" || len(cs.PreviousTrade) == 0 {
			if bearishHammer || bullishHammer || bullishMaru || isDozi || (isBull && !invertedHammer) {
				if (shortTrend == "decline" && shortTrendCount >= 2) || (bearTrendCount >= 2 || bearCount >= 2) || lhePattern >= 3 {
					log.Printf("BUY BajajFinance %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
					log.Printf("Previous Trade: %v :: Bearish Hammer: %v :: bullishHammer: %v :: isBull: %v :: BullishMaru:: %v :: isDozi: %v", cs.PreviousTrade, bearishHammer, bullishHammer, isBull, bullishMaru, isDozi)
					log.Printf("shortTrend: %v :: shortTrendCount: %v :: bearTrendCount: %v :: bearCount: %v :: lhePattern:: %v", shortTrend, shortTrendCount, bearTrendCount, bearCount, lhePattern)
					msg := fmt.Sprintf("BUY CALL ::  %s - %s - %s\nPrevious Day Trend: %s \nPrevious Day Change: %s \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v  \nMESSAGE : %s \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, "Ensure that market is moving up", separation)
					alerts.SendAlerts(msg, alerts.BuyStockChannel)
				}

			}

			if len(cs.PreviousTrade) == 0 {
				if invertedHammer || bearishMaru || isDozi || isBear {
					if (shortTrend == "rally" && shortTrendCount >= 2) || (bullTrendCount >= 2 || bullCount >= 2) || hhePattern >= 4 {
						log.Printf("SHORT SELL BajajFinance %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
						log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
						log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
						msg := fmt.Sprintf("SHORT SELL CALL :: %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v \nMESSAGE: %s  \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, "Ensure that market is falling down", separation)
						alerts.SendAlerts(msg, alerts.ShortSellStocksChannel)
					}

				}

			}

		} else if cs.PreviousTrade == "BOUGHT" {
			if isBear || bearishMaru || isDozi {
				if (shortTrend == "rally" && shortTrendCount >= 2) || (bullTrendCount >= 2 || bullCount >= 2) {
					log.Printf("SELL BajajFinance %s - %s - %s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange)
					log.Printf("Previous Trade: %v :: isBear: %v :: bearishMaru:  %v :: isDozi: %v", cs.PreviousTrade, isBear, bearishMaru, isDozi)
					log.Printf("shortTrend: %v :: shortTrendCount: %v :: bullTrendCount: %v :: bullCount: %v :: hhePattern:: %v", shortTrend, shortTrendCount, bullTrendCount, bullCount, hhePattern)
					msg := fmt.Sprintf("SELL CALL %s - %s - %s \nPrevious Day Trend: %s \nPrevious Day Change: %s  \nBuyQuanity: %v \nSellQuantity: %v \nChange: %v  \n%s", cs.Instrument.Name, cs.Instrument.Symbol, cs.Instrument.Exchange, prevDayTrend, prevDayChange, bq, sq, qchange, separation)
					alerts.SendAlerts(msg, alerts.SellStockChannel)
				}

			}
		}

	}

	log.Println("-----------------------------------------------------------------------------------------------------------------------")

	return "", nil

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
		change = fmt.Sprintf("+%.2f%%", ((close-open)/close)*100) //TODO: Send data only if change is greater than 1.5%
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

//GetOHLC return open high low and close for the stock
func (cs CandleStick) GetOHLC() (float64, float64, float64, error) {
	cs.Mux.Lock()
	defer cs.Mux.Unlock()
	time.Sleep(1500 * time.Millisecond) //wait to avoid to many request error
	token := cs.Instrument.Token
	ohlc, err := cs.KC.GetOHLC(token)
	if err != nil {
		log.Println("Error finding OHLC : ", err)
		return 0, 0, 0, err
	}

	return ohlc[token].OHLC.Open, ohlc[token].OHLC.High, ohlc[token].OHLC.Low, nil

}

//GetTradeQuantity returns total buy and sell trade and percentage change
func (cs CandleStick) GetTradeQuantity() (float64, float64, string, error) {
	cs.Mux.Lock()
	defer cs.Mux.Unlock()
	time.Sleep(1500 * time.Millisecond)
	quote, err := cs.KC.GetQuote(cs.Instrument.Token)
	if err != nil {
		log.Println("Error finding OHLC : ", err)
		return 0, 0, "", err
	}

	bq := float64(quote[cs.Instrument.Token].BuyQuantity)
	sq := float64(quote[cs.Instrument.Token].SellQuantity)

	var qChange string
	if bq > sq {

		qChange = fmt.Sprintf("+%.2f%%", ((bq-sq)/bq)*100)
	} else {
		qChange = fmt.Sprintf("-%.2f%%", ((sq-bq)/sq)*100)

	}

	return bq, sq, qChange, nil

}

//UnSubcribe a stock to stop getting ticks
func (cs CandleStick) UnSubcribe() error {
	var us []uint32
	us = append(us, getUnit32(cs.Instrument.Token))
	ticker := kiteticker.New(apiKey, accessToken)
	return ticker.Unsubscribe(us)
}

//SendProfitAlerts sends alerts about profit and loss for a particular postion.
func (cs CandleStick) SendProfitAlerts() error {
	var pl float64
	time.Sleep(2000 * time.Millisecond) //wait to avoid to many request error
	pos, err := cs.KC.GetPositions()
	if err != nil {
		return err
	}
	//log.Printf("Positions : %+v\n", pos)

	for _, pos := range pos.Net {
		if pos.InstrumentToken == getUnit32(cs.Instrument.Token) {
			pl = pos.Unrealised
		}

	}
	if pl != 0 {
		if pl < 0 {
			msg := fmt.Sprintf("PROFIT DROPPED TO NEGATIVE: \nInstrument: %s \nProfit-Loss: %.2f", cs.Instrument.Symbol, pl)
			alerts.SendAlerts(msg, alerts.TradeChannel)
		} else if pl > 1000 {
			msg := fmt.Sprintf("PROFIT above 1000 \nInstrument: %s \nProfit-Loss: %.2f ", cs.Instrument.Symbol, pl)
			alerts.SendAlerts(msg, alerts.TradeChannel)
		} else if pl > 2000 {
			msg := fmt.Sprintf("PROFIT above 2000 \nInstrument: %s  \nProfit-Loss: %.2f", cs.Instrument.Symbol, pl)
			alerts.SendAlerts(msg, alerts.TradeChannel)
		}
	} else {
		msg := fmt.Sprintf("Instrument %s does not have any open positions", cs.Instrument.Symbol)
		log.Println(msg)
	}

	return nil
}
