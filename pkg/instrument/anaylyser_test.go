package instrument

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var uptrendData1 = [][]float64{
	{50, 60, 45, 40},
	{45, 55, 40, 35},
	{40, 50, 35, 30},
	{35, 45, 30, 25},
	{30, 40, 25, 20},
	{25, 35, 20, 15},
}

var uptrendData2 = [][]float64{
	{50, 60, 45, 40},
	{45, 55, 40, 35},
	{40, 50, 35, 30},
	{35, 45, 30, 25},
	{30, 40, 25, 20},
	{40, 50, 35, 30},
	{25, 35, 20, 15},
}

var downtrendData1 = [][]float64{
	{25, 35, 20, 15},
	{30, 40, 25, 20},
	{35, 45, 30, 25},
	{40, 50, 35, 30},
	{45, 55, 40, 35},
	{50, 60, 45, 40},
}

var downTrendData2 = [][]float64{
	{25, 35, 20, 15},
	{30, 40, 25, 20},
	{35, 45, 30, 25},
	{40, 50, 35, 30},
	{45, 55, 40, 35},
	{35, 45, 30, 25},
	{50, 60, 45, 40},
}
var downtrendData = [][]float64{}

func getTrendList(data [][]float64) []TradeDetails {
	var tdList []TradeDetails
	for _, d := range data {
		td := &TradeDetails{
			Open:  d[0],
			High:  d[1],
			Close: d[2],
			Low:   d[3],
		}
		tdList = append(tdList, *td)
	}
	return tdList

}
func getTestData(open, high, close, low float64) *TradeDetails {
	td := &TradeDetails{
		Open:  open,
		High:  high,
		Close: close,
		Low:   low,
	}

	return td

}
func TestIsBullish(t *testing.T) {
	var tdList []TradeDetails
	td := getTestData(10, 20, 15, 5)
	tdList = append(tdList, *td)
	res, count := isBullish(tdList)
	assert.True(t, res)
	assert.Equal(t, 1, count)

	res, count = isBearish(tdList)
	assert.False(t, res)
	assert.Equal(t, 0, count)

}

func TestBullishCount(t *testing.T) {
	var tdList []TradeDetails
	td1 := getTestData(10, 20, 15, 5)
	tdList = append(tdList, *td1)
	td2 := getTestData(12, 22, 17, 7)
	tdList = append(tdList, *td2)
	td3 := getTestData(14, 24, 19, 9)
	tdList = append(tdList, *td3)
	res, count := isBullish(tdList)
	td4 := getTestData(20, 25, 15, 5)
	tdList = append(tdList, *td4)
	assert.True(t, res)
	assert.Equal(t, 3, count)

	res, count = isBearish(tdList)
	assert.False(t, res)
	assert.Equal(t, 0, count)
}

func TestIsBearish(t *testing.T) {
	var tdList []TradeDetails
	td := getTestData(20, 25, 15, 5)
	tdList = append(tdList, *td)
	res, count := isBearish(tdList)
	assert.True(t, res)
	assert.Equal(t, 1, count)

	res, count = isBullish(tdList)
	assert.False(t, res)
	assert.Equal(t, 0, count)

}

func TestIsBearishCount(t *testing.T) {
	var tdList []TradeDetails
	td1 := getTestData(20, 25, 15, 5)
	tdList = append(tdList, *td1)
	td2 := getTestData(18, 23, 13, 3)
	tdList = append(tdList, *td2)
	td3 := getTestData(16, 21, 11, 2)
	tdList = append(tdList, *td3)
	td4 := getTestData(10, 20, 15, 5)
	tdList = append(tdList, *td4)
	res, count := isBearish(tdList)
	assert.True(t, res)
	assert.Equal(t, 3, count)

	res, count = isBullish(tdList)
	assert.False(t, res)
	assert.Equal(t, 0, count)
}

func TestIsBullishMarubuzo(t *testing.T) {
	//Positive Test
	td := getTestData(10, 20, 20, 10)
	res := isBullishMarubuzo(*td)
	assert.True(t, res)

	//Bullish Body > 80% Test
	td2 := getTestData(10, 22, 20, 10)
	res2 := isBullishMarubuzo(*td2)
	assert.True(t, res2)

	td5 := getTestData(10, 21, 20, 9)
	res5 := isBullishMarubuzo(*td5)
	assert.True(t, res5)

	//Negative Test
	td3 := getTestData(15, 20, 20, 10)
	res3 := isBullishMarubuzo(*td3)
	assert.False(t, res3)

	td4 := getTestData(10, 20, 10, 10)
	res4 := isBullishMarubuzo(*td4)
	assert.False(t, res4)

}

func TestIsBearishMarubuzo(t *testing.T) {
	td := getTestData(20, 20, 10, 10)
	res := isBearishMarubuzo(*td)
	assert.True(t, res)

	//Bearish body greater than 80%
	td2 := getTestData(20, 21, 10, 10)
	res2 := isBearishMarubuzo(*td2)
	assert.True(t, res2)

	td3 := getTestData(20, 20, 10, 9)
	res3 := isBearishMarubuzo(*td3)
	assert.True(t, res3)

	td4 := getTestData(20, 21, 10, 11)
	res4 := isBearishMarubuzo(*td4)
	assert.True(t, res4)

	td5 := getTestData(20, 25, 10, 11)
	res5 := isBearishMarubuzo(*td5)
	assert.False(t, res5)
}

func TestIsDozi(t *testing.T) {
	td := getTestData(10, 20, 10, 10)
	res := isDozi(*td)
	assert.True(t, res)

	td2 := getTestData(10, 10, 10, 5)
	res2 := isDozi(*td2)
	assert.True(t, res2)

	td3 := getTestData(10, 12, 10, 5)
	res3 := isDozi(*td3)
	assert.True(t, res3)

}

func TestIsInvertedHammer(t *testing.T) {
	// Positive hammer
	td := getTestData(5, 25, 10, 2)
	res := isInvertedHammer(*td)
	assert.True(t, res)

	td2 := getTestData(5, 19, 10, 2)
	res2 := isInvertedHammer(*td2)
	assert.False(t, res2)

	// Negative Hammer
	td3 := getTestData(10, 25, 5, 2)
	res3 := isInvertedHammer(*td3)
	assert.True(t, res3)

	td4 := getTestData(10, 19, 5, 2)
	res4 := isInvertedHammer(*td4)
	assert.False(t, res4)
}

func TestIsBullishHammer(t *testing.T) {
	// Bullish hammer
	td := getTestData(20, 27, 25, 9)
	res := isBullishHammer(*td)
	assert.True(t, res)

	td2 := getTestData(20, 27, 25, 10)
	res2 := isBullishHammer(*td2)
	assert.True(t, res2)

	td3 := getTestData(20, 27, 25, 9)
	res3 := isBullishHammer(*td3)
	assert.True(t, res3)

	td4 := getTestData(20, 27, 25, 11)
	res4 := isBullishHammer(*td4)
	assert.False(t, res4)

}

func TestIsBearishHammer(t *testing.T) {
	// Bearish Hammer
	td3 := getTestData(25, 27, 20, 9)
	res3 := isBearishHammer(*td3)
	assert.True(t, res3)

	td4 := getTestData(25, 27, 20, 10)
	res4 := isBearishHammer(*td4)
	assert.True(t, res4)

	td5 := getTestData(25, 27, 20, 8)
	res5 := isBearishHammer(*td5)
	assert.True(t, res5)

	td6 := getTestData(25, 27, 20, 11)
	res6 := isBearishHammer(*td6)
	assert.False(t, res6)
}

func TestUpTrend(t *testing.T) {
	td := getTrendList(uptrendData1)
	trend, count := getShortTermTrend(td)
	assert.Equal(t, "rally", trend)
	assert.Equal(t, 5, count)

	td2 := getTrendList(uptrendData2)

	trend2, count2 := getShortTermTrend(td2)
	assert.Equal(t, "rally", trend2)
	assert.Equal(t, 4, count2)
}

func TestDownTrend(t *testing.T) {
	td := getTrendList(downtrendData1)
	trend, count := getShortTermTrend(td)
	assert.Equal(t, "decline", trend)
	assert.Equal(t, 5, count)

	td2 := getTrendList(downTrendData2)

	trend2, count2 := getShortTermTrend(td2)
	assert.Equal(t, "decline", trend2)
	assert.Equal(t, 4, count2)
}

func TestTime(t *testing.T) {
	getActualMarketOpenTime("2019-06-09")
}
