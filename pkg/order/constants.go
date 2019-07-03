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
		Symbol:       "INDIGO",
		Token:        "2865921",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1552,
		SellStopLoss: 1565,
		SellTarget1:  1542,
		SellTarget2:  1532,
		Quantity:     200,
		BuyPrice:     1565,
		BuyStopLoss:  1558,
		BuyTarget1:   1575,
		BuyTarget2:   1585,
		MaxLoss:      -700,
	},

	{
		Symbol:       "INDUSINDBK",
		Token:        "1346049",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1410,
		SellStopLoss: 1417,
		SellTarget1:  1400,
		SellTarget2:  1390,
		Quantity:     200,
		BuyPrice:     1430,
		BuyStopLoss:  1423,
		BuyTarget1:   1440,
		BuyTarget2:   1450,
		MaxLoss:      -700,
	},

	// {
	// 	Symbol:       "IBULHSGFIN",
	// 	Token:        "7712001",
	// 	Exchange:     "NSE",
	// 	Variety:      "bo",
	// 	Product:      "MIS",
	// 	OrderType:    "LIMIT",
	// 	Validity:     "DAY",
	// 	SellPrice:    630,
	// 	SellStopLoss: 637,
	// 	SellTarget1:  620,
	// 	SellTarget2:  610,
	// 	Quantity:     200,
	// 	BuyPrice:     649,
	// 	BuyStopLoss:  642,
	// 	BuyTarget1:   659,
	// 	BuyTarget2:   669,
	// 	MaxLoss:      -700,
	// },

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
}
