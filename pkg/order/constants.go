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
}

//BuyOrders a list of buy ordres
var BuyOrders = []BOInput{

	// {
	// 	Symbol:       "ADANIPORTS",
	// 	Token:        "3861249",
	// 	Exchange:     "NSE",
	// 	Variety:      "bo",
	// 	Product:      "MIS",
	// 	OrderType:    "LIMIT",
	// 	Validity:     "DAY",
	// 	SellPrice:    0,
	// 	SellStopLoss: 0,
	// 	SellTarget1:  0,
	// 	SellTarget2:  0,
	// 	Quantity:     500,
	// 	BuyPrice:     414,
	// 	BuyStopLoss:  413,
	// 	BuyTarget1:   417,
	// 	BuyTarget2:   417,
	// },

	{
		Symbol:       "UBL",
		Token:        "4278529",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1352,
		SellStopLoss: 1362,
		SellTarget1:  1342,
		SellTarget2:  1332,
		Quantity:     200,
		BuyPrice:     1365,
		BuyStopLoss:  1355,
		BuyTarget1:   1375,
		BuyTarget2:   1380,
	},

	{
		Symbol:       "JUSTDIAL",
		Token:        "7670273",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    758,
		SellStopLoss: 768,
		SellTarget1:  748,
		SellTarget2:  738,
		Quantity:     200,
		BuyPrice:     769,
		BuyStopLoss:  759,
		BuyTarget1:   779,
		BuyTarget2:   785,
	},
	{
		Symbol:       "DIVISLAB",
		Token:        "2800641",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		SellPrice:    1600,
		SellStopLoss: 1610,
		SellTarget1:  1590,
		SellTarget2:  1590,
		Quantity:     200,
		BuyPrice:     1614,
		BuyStopLoss:  1604,
		BuyTarget1:   1624,
		BuyTarget2:   1624,
	},

	// {
	// 	Symbol:       "ESCORTS",
	// 	Token:        "245249",
	// 	Exchange:     "NSE",
	// 	Variety:      "bo",
	// 	Product:      "MIS",
	// 	OrderType:    "LIMIT",
	// 	Validity:     "DAY",
	// 	SellPrice:    0,
	// 	SellStopLoss: 0,
	// 	SellTarget1:  0,
	// 	SellTarget2:  0,
	// 	Quantity:     500,
	// 	BuyPrice:     574,
	// 	BuyStopLoss:  572,
	// 	BuyTarget1:   580,
	// 	BuyTarget2:   580,
	// },
}
