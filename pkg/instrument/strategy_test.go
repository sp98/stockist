package instrument

import (
	"log"
	"testing"
)

var testData1 = [][]float64{
	//Open, High, Low, Close, AverageTradedPrice (in reverse order. )
	{499.7, 499.7, 496.6, 498.25, 497.79},
	{501, 501, 500, 500, 500},
	{501, 501, 501, 501, 0},
	{494.1, 502.8, 489.85, 501, 0},
}

var testData2 = [][]float64{
	//Open, High, Low, Close, AverageTradedPrice (in reverse order. )
	{378.2, 380.15, 377.05, 379.1, 378.89},
	{385.2, 385.2, 377, 377, 377},
	{385.2, 385.2, 385.2, 385.2, 0},
	{378.1, 386.45, 377, 385.2, 0},
}

var testData3 = [][]float64{
	//Open, High, Low, Close, AverageTradedPrice (in reverse order. )
	{286.6, 287.8, 285.5, 286.9, 286.26},
	{292.65, 292.65, 286, 286, 286},
	{292.65, 292.65, 292.65, 292.65, 0},
	{297, 299.6, 289.25, 292.65, 0},
}

var testData4 = [][]float64{
	//Open, High, Low, Close, AverageTradedPrice (in reverse order. )
	{1555.25, 1559.55, 1537.8, 1545.9, 1545.81},
	{1582.25, 1582.25, 1559.9, 1559.9, 1559.9},
	{1582.25, 1582.25, 1582.25, 1582.25, 0},
	{1622, 1623.95, 1564, 1582.25, 0},
}

var testData5 = [][]float64{
	//Open, High, Low, Close, AverageTradedPrice (in reverse order. )
	// {179.9, 180, 178.05, 178.8, 178.84},
	{181.35, 181.35, 179.95, 179.95, 179.99},
	{181.35, 181.35, 181.35, 181.35, 0},
	{176.4, 183.25, 173.35, 181.35, 0},
}

var testData6 = [][]float64{
	//Open, High, Low, Close, AverageTradedPrice (in reverse order. )
	{179.9, 180, 178.05, 178.8, 0},
	{181.35, 181.35, 179.95, 179.95, 0},
	{181.35, 181.35, 181.35, 181.35, 0},
	{176.4, 183.25, 173.35, 181.35, 0},
}

var dataBearishMaruAfterRally = [][]float64{
	//open High Low Close AverageTradedPrice
	{40, 40, 30, 30, 0}, //BullishHammer
	{35, 55, 25, 45, 0},
	{30, 50, 20, 40, 0},
	{25, 45, 15, 35, 0},
	{20, 40, 10, 30, 0},
}

var dataBearishInvertedHammerAfterRally = [][]float64{
	//open High Low Close AverageTradedPrice
	{25, 40, 23, 20, 0}, //InvertedHammer
	{35, 55, 25, 45, 0},
	{30, 50, 20, 40, 0},
	{25, 45, 15, 35, 0},
	{20, 40, 10, 30, 0},
}

var dataBullishInvertedHammerAfterRally = [][]float64{
	//open High Low Close AverageTradedPrice
	{20, 40, 23, 25, 0}, //BullishInvertedHammer
	{35, 55, 25, 45, 0},
	{30, 50, 20, 40, 0},
	{25, 45, 15, 35, 0},
	{20, 40, 10, 30, 0},
}

// var downtrendSensexData = [][]float64{
// 	//open, high, close, low
// 	{25, 35, 15, 20, 0},
// 	{30, 40, 20, 25, 0},
// 	{35, 45, 25, 30, 0},
// 	{40, 50, 30, 35, 0},
// 	{45, 55, 35, 40, 0},
// 	{50, 60, 40, 45, 0},
// }

func getTesCandleStick(data [][]float64) CandleStick {
	var csList []CandleStickList
	inst := &Instrument{
		Name:     "ACC",
		Exchange: "NSE",
		Symbol:   "ACC",
		Interval: "5m",
		Token:    "5633",
	}
	for _, d := range data {
		td := &CandleStickList{
			Open:               d[0],
			High:               d[1],
			Low:                d[2],
			Close:              d[3],
			AverageTradedPrice: d[4],
		}
		csList = append(csList, *td)
	}
	cs := &CandleStick{
		KC:         getConnection(),
		Instrument: *inst,
		Details:    csList,
	}
	return *cs

}

func TestOpeningTrend(t *testing.T) {
	// cs := getTesCandleStick(testData1)
	// cs.OpeningTrend()

	// cs3 := getTesCandleStick(testData3)
	// cs3.OpeningTrend()

	// cs2 := getTesCandleStick(testData2)
	// cs2.OpeningTrend()
	// cs4 := getTesCandleStick(testData4)
	// cs4.OpeningTrend()
	cs5 := getTesCandleStick(testData2)
	cs5.OpeningTrend()

}

func TestAnalyseSensex(t *testing.T) {
	// cs := getTesCandleStick(dataBearishMaruAfterRally)
	// cs.AnalyseSensex()

	//cs2 := getTesCandleStick(downtrendSensexData)
	//cs2.AnalyseSensex()

	cs3 := getTesCandleStick(dataBearishInvertedHammerAfterRally)
	cs3.AnalyseSensex()

	cs4 := getTesCandleStick(dataBullishInvertedHammerAfterRally)
	cs4.AnalyseSensex()

}

func TestPreviousDayTrend(t *testing.T) {
	trend, change := getPreviousDayTrend(4, 2)
	t.Error(trend)
	t.Error(change)
}

func TestUnSubscribe(t *testing.T) {
	cs := getTesCandleStick(dataBearishInvertedHammerAfterRally)
	log.Printf("CS %+v", cs)
	q, err := cs.KC.GetQuote(cs.Instrument.Token)
	if err != nil {
		t.Error(q)
	}
	t.Error(q[cs.Instrument.Token].BuyQuantity)
	t.Error(q[cs.Instrument.Token].SellQuantity)
}
