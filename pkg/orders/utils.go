package orders

func isBullish(tradeDetails []TradeDetails) (bool, int) {

	isBull := true
	lastCandleStick := tradeDetails[0]
	if lastCandleStick.Open >= lastCandleStick.Close {
		isBull = false
		return isBull, 0
	}

	//log.Println("Check bull count!")
	var count int
	for _, td := range tradeDetails {
		//log.Printf("%d :: Open- %f  Close- %f", i, td.Open, td.Close)
		if td.Open > td.Close {
			break
		}
		count = count + 1
	}

	return isBull, count

}

func isBearish(tradeDetails []TradeDetails) (bool, int) {

	isBear := true
	lastCandleStick := tradeDetails[0]
	if lastCandleStick.Open < lastCandleStick.Close {
		isBear = false
		return isBear, 0
	}
	var trendCount int
	for _, td := range tradeDetails {
		if td.Open <= td.Close {
			break
		}
		trendCount = trendCount + 1
	}

	return isBear, trendCount

}

func isBullishMarubuzo(td TradeDetails) bool {
	if td.Open < td.Close {
		//println(((td.Close - td.Open) / (td.High - td.Low)))
		if td.Open == td.Low && td.Close == td.High {
			return true
		}

		if (((td.Close - td.Open) / (td.High - td.Low)) * 100) > 80 {
			return true
		}
	}

	return false

}

func isBearishMarubuzo(td TradeDetails) bool {

	if td.Open > td.Close {
		//println(((td.Open - td.Close) / (td.High - td.Low)))
		if td.Open == td.High && td.Close == td.Low {
			return true
		}

		if (((td.Open - td.Close) / (td.High - td.Low)) * 100) > 80 {
			return true
		}
	}

	return false

}

func isDozi(td TradeDetails) bool {
	if td.Open == td.Close && (td.High != td.Open || td.Low != td.Open) {
		return true
	}
	// else if (td.Open == td.Close) && (td.Close == td.High) && (td.Low != td.Open) {
	// 	return true
	// } else if (td.Open == td.Close) && (td.Close == td.Low) && (td.Low != td.High) {
	// 	return true
	// }
	return false

}

func isInvertedHammer(td TradeDetails) bool {
	if td.Open < td.Close {
		if (2*(td.Close-td.Open) < (td.High - td.Close)) && ((td.Open - td.Low) < (td.Close - td.Open)) {
			return true
		}
	} else if td.Open > td.Close {
		if (2*(td.Open-td.Close) < (td.High - td.Open)) && ((td.Close - td.Low) < (td.Open - td.Close)) {
			return true
		}
	}

	return false
}

func isBullishHammer(td TradeDetails) bool {
	if td.Open < td.Close {
		if ((td.Open - td.Low) >= 2*(td.Close-td.Open)) && ((td.High - td.Close) < (td.Close - td.Open)) {
			return true
		}
	}

	return false
}

func isBearishHammer(td TradeDetails) bool {
	if td.Open > td.Close {
		if ((td.Close - td.Low) >= 2*(td.Open-td.Close)) && ((td.High - td.Open) < (td.Open - td.Close)) {
			return true
		}
	}
	return false
}

//isLowerHighsEngulfingPatter checks for patter where lower highs are made but lows may be lower or higher (making the previous pattern engulfuing the previous one)
func lowerHighsEngulfingPatternCount(td []TradeDetails) int {
	count := 0
	for i := 0; i < len(td)-1; i++ {
		if td[i].High < td[i+1].High && ((td[i].Low < td[i+1].Low) || (td[i].Low > td[i+1].Low)) {
			count = count + 1
			continue
		}
		return count
	}
	return count
}

func higherLowsEngulfingPatternCount(td []TradeDetails) int {
	count := 0
	for i := 0; i < len(td)-1; i++ {
		if td[i].Low > td[i+1].Low && ((td[i].High > td[i+1].High) || (td[i].High <= td[i+1].High)) {
			count = count + 1
			continue
		}
		return count
	}
	return count
}
