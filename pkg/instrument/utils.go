package instrument

func isBullish(csDetails []CandleStickDetails) (bool, int) {

	isBull := true
	lastCandleStick := csDetails[0]
	if lastCandleStick.Open >= lastCandleStick.Close {
		isBull = false
		return isBull, 0
	}

	//log.Println("Check bull count!")
	var count int
	for _, cs := range csDetails {
		//log.Printf("%d :: Open- %f  Close- %f", i, cs.Open, cs.Close)
		if cs.Open > cs.Close {
			break
		}
		count = count + 1
	}

	return isBull, count

}

func isBearish(csDetails []CandleStickDetails) (bool, int) {

	isBear := true
	lastCandleStick := csDetails[0]
	if lastCandleStick.Open < lastCandleStick.Close {
		isBear = false
		return isBear, 0
	}
	var trendCount int
	for _, cs := range csDetails {
		if cs.Open <= cs.Close {
			break
		}
		trendCount = trendCount + 1
	}

	return isBear, trendCount

}

func isBullishMarubuzo(csDetails CandleStickDetails) bool {
	if csDetails.Open < csDetails.Close {
		//println(((cs.Close - cs.Open) / (cs.High - cs.Low)))
		if csDetails.Open == csDetails.Low && csDetails.Close == csDetails.High {
			return true
		}

		if (((csDetails.Close - csDetails.Open) / (csDetails.High - csDetails.Low)) * 100) > 80 {
			return true
		}
	}

	return false

}

func isBearishMarubuzo(cs CandleStickDetails) bool {

	if cs.Open > cs.Close {
		//println(((cs.Open - cs.Close) / (cs.High - cs.Low)))
		if cs.Open == cs.High && cs.Close == cs.Low {
			return true
		}

		if (((cs.Open - cs.Close) / (cs.High - cs.Low)) * 100) > 80 {
			return true
		}
	}

	return false

}

func isDozi(cs CandleStickDetails) bool {
	if cs.Open == cs.Close && (cs.High != cs.Open || cs.Low != cs.Open) {
		return true
	}
	// else if (cs.Open == cs.Close) && (cs.Close == cs.High) && (cs.Low != cs.Open) {
	// 	return true
	// } else if (cs.Open == cs.Close) && (cs.Close == cs.Low) && (cs.Low != cs.High) {
	// 	return true
	// }
	return false

}

func isInvertedHammer(cs CandleStickDetails) bool {
	if cs.Open < cs.Close {
		if (2*(cs.Close-cs.Open) < (cs.High - cs.Close)) && ((cs.Open - cs.Low) < (cs.Close - cs.Open)) {
			return true
		}
	} else if cs.Open > cs.Close {
		if (2*(cs.Open-cs.Close) < (cs.High - cs.Open)) && ((cs.Close - cs.Low) < (cs.Open - cs.Close)) {
			return true
		}
	}

	return false
}

func isBullishHammer(cs CandleStickDetails) bool {
	if cs.Open < cs.Close {
		if ((cs.Open - cs.Low) >= 2*(cs.Close-cs.Open)) && ((cs.High - cs.Close) < (cs.Close - cs.Open)) {
			return true
		}
	}

	return false
}

func isBearishHammer(cs CandleStickDetails) bool {
	if cs.Open > cs.Close {
		if ((cs.Close - cs.Low) >= 2*(cs.Open-cs.Close)) && ((cs.High - cs.Open) < (cs.Open - cs.Close)) {
			return true
		}
	}
	return false
}

//isLowerHighsEngulfingPatter checks for patter where lower highs are made but lows may be lower or higher (making the previous pattern engulfuing the previous one)
func lowerHighsEngulfingPatternCount(cs []CandleStickDetails) int {
	count := 0
	for i := 0; i < len(cs)-1; i++ {
		if cs[i].High < cs[i+1].High && ((cs[i].Low < cs[i+1].Low) || (cs[i].Low > cs[i+1].Low)) {
			count = count + 1
			continue
		}
		return count
	}
	return count
}

func higherLowsEngulfingPatternCount(cs []CandleStickDetails) int {
	count := 0
	for i := 0; i < len(cs)-1; i++ {
		if cs[i].Low > cs[i+1].Low && ((cs[i].High > cs[i+1].High) || (cs[i].High <= cs[i+1].High)) {
			count = count + 1
			continue
		}
		return count
	}
	return count
}
