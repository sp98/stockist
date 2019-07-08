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

	{
		Symbol:       "INDUSINDBK",
		Token:        "1346049",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		BuyPrice:     1525,
		BuyStopLoss:  10,
		BuyTarget1:   10,
		BuyTarget2:   20,
		SellPrice:    1505,
		SellStopLoss: 10,
		SellTarget1:  10,
		SellTarget2:  20,
		Quantity:     100,
		MaxLoss:      -800,
	},
	{
		Symbol:       "IBULHSGFIN",
		Token:        "7712001",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		BuyPrice:     740,
		BuyStopLoss:  10,
		BuyTarget1:   10,
		BuyTarget2:   20,
		SellPrice:    728,
		SellStopLoss: 10,
		SellTarget1:  10,
		SellTarget2:  20,
		Quantity:     100,
		MaxLoss:      -800,
	},
	{
		Symbol:       "RELIANCE",
		Token:        "738561",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		BuyPrice:     1266,
		BuyStopLoss:  10,
		BuyTarget1:   10,
		BuyTarget2:   20,
		SellPrice:    1254,
		SellStopLoss: 10,
		SellTarget1:  10,
		SellTarget2:  20,
		Quantity:     100,
		MaxLoss:      -800,
	},
	{
		Symbol:       "SIEMENS",
		Token:        "806401",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		BuyPrice:     1282,
		BuyStopLoss:  10,
		BuyTarget1:   10,
		BuyTarget2:   20,
		SellPrice:    1263,
		SellStopLoss: 10,
		SellTarget1:  10,
		SellTarget2:  20,
		Quantity:     100,
		MaxLoss:      -800,
	},
	{
		Symbol:       "INDIGO",
		Token:        "2865921",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		BuyPrice:     1590,
		BuyStopLoss:  10,
		BuyTarget1:   10,
		BuyTarget2:   20,
		SellPrice:    1558,
		SellStopLoss: 10,
		SellTarget1:  10,
		SellTarget2:  20,
		Quantity:     100,
		MaxLoss:      -800,
	},
	{
		Symbol:       "TITAN",
		Token:        "897537",
		Exchange:     "NSE",
		Variety:      "bo",
		Product:      "MIS",
		OrderType:    "LIMIT",
		Validity:     "DAY",
		BuyPrice:     1280,
		BuyStopLoss:  10,
		BuyTarget1:   10,
		BuyTarget2:   20,
		SellPrice:    1254,
		SellStopLoss: 10,
		SellTarget1:  10,
		SellTarget2:  20,
		Quantity:     100,
		MaxLoss:      -800,
	},
}
