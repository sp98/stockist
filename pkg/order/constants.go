package order

//BOInput is the input struct for Buy Orders
type BOInput struct {
	Symbol          string
	Token           string
	Exchange        string
	Product         string
	Price           float64
	StopLoss        float64
	Target1         float64
	Target2         float64
	TransactionType string
	OrderType       string
	Variety         string
	Quantity        int
	Validity        string
}

//SOInput is the input struct for Sell Orders
type SOInput struct {
	Symbol          string
	Token           string
	Exchange        string
	Product         string
	SellPrice       float64
	SellStopLoss    float64
	SellTarget1     float64
	SellTarget2     float64
	TransactionType string
	OrderType       string
	Variety         string
	Quantity        int
	Validity        string
}

//BuyOrders a list of buy ordres
var BuyOrders = []BOInput{
	// {
	// 	Symbol:          "ADANIPORTS",
	// 	Token:           "3861249",
	// 	Exchange:        "NSE",
	// 	Variety:         "bo",
	// 	Product:         "MIS",
	// 	OrderType:       "LIMIT",
	// 	TransactionType: "BUY",
	// 	Validity:        "DAY",
	// 	Price:           409,
	// 	StopLoss:        405,
	// 	Target1:         418,
	// 	Target2:         422,
	// },
	{
		Symbol:          "BALKRISIND",
		Token:           "85761",
		Exchange:        "NSE",
		Variety:         "bo",
		Product:         "MIS",
		OrderType:       "LIMIT",
		TransactionType: "BUY",
		Validity:        "DAY",
		Price:           765,
		StopLoss:        755,
		Target1:         775,
		Target2:         785,
		Quantity:        1,
	},

	{
		Symbol:          "BALKRISIND",
		Token:           "85761",
		Exchange:        "NSE",
		Variety:         "bo",
		Product:         "MIS",
		OrderType:       "LIMIT",
		TransactionType: "SELL",
		Validity:        "DAY",
		Price:           753,
		StopLoss:        763,
		Target1:         743,
		Target2:         733,
		Quantity:        1,
	},
}
