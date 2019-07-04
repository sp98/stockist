package order

//BOInput is the input struct for Buy Orders
type BOInput struct {
	Symbol       string
	Token        string
	Exchange     string
	Product      string
	BuyPrice     float64
	BuyStopLoss  float64
	BuyTarget1   float64
	BuyTarget2   float64
	SellPrice    float64
	SellStopLoss float64
	SellTarget1  float64
	SellTarget2  float64
	OrderType    string
	Variety      string
	Quantity     int
	Validity     string
	MaxLoss      float64
}

//BuyOrders a list of buy ordres
var BuyOrders = []BOInput{

	// {
	// 	Symbol:       "BHARATFIN",
	// 	Token:        "4995329",
	// 	Exchange:     "NSE",
	// 	Variety:      "bo",
	// 	Product:      "MIS",
	// 	OrderType:    "LIMIT",
	// 	Validity:     "DAY",
	// 	SellPrice:    894,
	// 	SellStopLoss: 901,
	// 	SellTarget1:  884,
	// 	SellTarget2:  874,
	// 	Quantity:     200,
	// 	BuyPrice:     907,
	// 	BuyStopLoss:  900,
	// 	BuyTarget1:   917,
	// 	BuyTarget2:   927,
	// 	MaxLoss:      -700,
	// },

	{
		Symbol:       "UBL",
		Token:        "4278529",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1376,
		SellStopLoss: 1388,
		SellTarget1:  1366,
		SellTarget2:  1356,
		Quantity:     200,
		BuyPrice:     1396,
		BuyStopLoss:  1388,
		BuyTarget1:   1406,
		BuyTarget2:   1416,
		MaxLoss:      -700,
	},

	{
		Symbol:       "INDIGO",
		Token:        "2865921",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1573,
		SellStopLoss: 1583,
		SellTarget1:  1563,
		SellTarget2:  1553,
		Quantity:     200,
		BuyPrice:     1585,
		BuyStopLoss:  1575,
		BuyTarget1:   1595,
		BuyTarget2:   1605,
		MaxLoss:      -700,
	},

	{
		Symbol:       "HINDUNILVR",
		Token:        "356865",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1776,
		SellStopLoss: 1788,
		SellTarget1:  1766,
		SellTarget2:  1756,
		Quantity:     200,
		BuyPrice:     1789,
		BuyStopLoss:  1781,
		BuyTarget1:   1799,
		BuyTarget2:   1809,
		MaxLoss:      -700,
	},

	{
		Symbol:       "IBULHSGFIN",
		Token:        "7712001",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    682,
		SellStopLoss: 690,
		SellTarget1:  672,
		SellTarget2:  662,
		Quantity:     200,
		BuyPrice:     696,
		BuyStopLoss:  686,
		BuyTarget1:   706,
		BuyTarget2:   716,
		MaxLoss:      -700,
	},

	// {
	// 	Symbol:       "DIVISLAB",
	// 	Token:        "2800641",
	// 	Exchange:     "NSE",
	// 	Variety:      "bo",
	// 	Product:      "MIS",
	// 	OrderType:    "LIMIT",
	// 	Validity:     "DAY",
	// 	SellPrice:    1600,
	// 	SellStopLoss: 1607,
	// 	SellTarget1:  1590,
	// 	SellTarget2:  1580,
	// 	Quantity:     200,
	// 	BuyPrice:     1621,
	// 	BuyStopLoss:  1614,
	// 	BuyTarget1:   1631,
	// 	BuyTarget2:   1641,
	// 	MaxLoss:      -700,
	// },

	{
		Symbol:       "HDFC",
		Token:        "340481",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    2272,
		SellStopLoss: 2280,
		SellTarget1:  2262,
		SellTarget2:  2252,
		Quantity:     200,
		BuyPrice:     2285,
		BuyStopLoss:  2277,
		BuyTarget1:   2295,
		BuyTarget2:   2305,
		MaxLoss:      -700,
	},
}
