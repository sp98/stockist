package order

import (
	"fmt"
	"log"
	"os"
	"strconv"

	kiteconnect "github.com/sp98/gokiteconnect"
)

var (
	apiKey              = os.Getenv("APIKEY")
	apiSecret           = os.Getenv("APISECRET")
	accessToken         = os.Getenv("ACCESSTOKEN")
	balance     float64 = 74000
)

func getConnection() *kiteconnect.Client {
	kc := kiteconnect.New(apiKey)
	kc.SetAccessToken(accessToken)
	return kc
}

//GetMargins returns the available margins
func (ord Order) GetMargins() {
	margins, err := ord.KC.GetUserMargins()
	if err != nil {
		fmt.Printf("Error getting margins: %v", err)
	}
	fmt.Printf("margins: %+v\n", margins)

}

//GetLastTradingPrice returns the last tarding price of the instruments
func (ord Order) GetLastTradingPrice() (float64, error) {
	res, err := ord.KC.GetLTP(ord.Token)
	if err != nil {
		log.Fatal("Error getting last traded price", err)
		return 0, err
	}
	return res[ord.Token].LastPrice, nil

}

//GetUnRealisedProfit returns unrealised profit and loss for the isntrument.
func (ord Order) GetUnRealisedProfit() (float64, error) {
	pos, _ := ord.KC.GetPositions()
	//log.Printf("Positions : %+v\n", pos)

	for _, pos := range pos.Net {
		if pos.InstrumentToken == getUnit32(ord.Token) {
			return pos.Unrealised, nil
		}

	}
	errMsg := fmt.Errorf("Error finding token %s in the Positions", ord.Token)
	log.Println(errMsg)
	return 0, errMsg
}

//GetOpenPrice returns the opening price of the stock
func (ord Order) GetOpenPrice() (float64, error) {
	ohlc, err := ord.KC.GetOHLC(ord.Token)
	if err != nil {
		log.Println("Error finding OHLC : ", err)
		return 0, err
	}
	// log.Printf("OHCL %+v", ohlc)
	return ohlc[ord.Token].OHLC.Open, nil

}

//GetClosePrice returns the opening price of the stock
func (ord Order) GetClosePrice() (float64, error) {
	ohlc, err := ord.KC.GetOHLC(ord.Token)
	if err != nil {
		log.Println("Error finding OHLC : ", err)
		return 0, err
	}
	// log.Printf("OHCL %+v", ohlc)
	return ohlc[ord.Token].OHLC.Close, nil

}

//ExecuteOrder executes an order
func (ord Order) ExecuteOrder() (*kiteconnect.OrderResponse, error) {
	log.Printf("Executing Order : %+v\n", ord)
	log.Printf("Order Params : %+v\n", ord.Params)
	ltp, _ := ord.GetLastTradingPrice() //Update price to Latest LTP
	ord.Params.Price = ltp
	res, err := ord.KC.PlaceOrder(ord.Variety, *ord.Params)
	if err != nil {
		log.Println("Error: Order execution failed :", err)
		return nil, err
	}

	log.Printf("Order Execution Response : %+v\n", res)

	return &res, nil

}

//GetParentOrderID returns parent Order for a given Order
func (ord Order) GetParentOrderID(id string) (string, error) {
	orders, err := ord.KC.GetOrders()
	if err != nil {
		log.Println("Error retrieving parent Order ID")
		return "", err
	}

	for _, order := range orders {
		if order.OrderID == id {
			log.Printf("Order found - %+v", order)
			return order.ParentOrderID, nil
		}
	}

	log.Println("Parent Order ID not found")

	return "", nil

}

//OrdersComplete checks whether all orders are complete or not
func (ord Order) OrdersComplete(ordResps []kiteconnect.OrderResponse) {
	//totalOrders := len(ordResps)
	keepLoop := true

	for keepLoop {
		orders, err := ord.KC.GetOrders()
		if err != nil {
			log.Println("Error retrieving parent Order ID")
			//return "", err
		}

		for _, order := range orders {
			for _, ordResp := range ordResps {
				if order.OrderID == ordResp.OrderID {
					//log.Printf("Order found - %+v", order)
					//return order.ParentOrderID, nil
				}
			}

		}

		log.Println("Parent Order ID not found")

		//return "", nil

	}

}

//getOrderStatus returns the status of the order.
func (ord Order) getOrderStatus(id string) (string, error) {
	var status string
	var message string
	orders, err := ord.KC.GetOrders()
	if err != nil {
		log.Println("Error retrieving parent Order ID")
		return "", err
	}

	for _, order := range orders {
		if order.OrderID == id {
			status = order.Status
			message = order.StatusMessage
		}
	}

	if status == "REJECTED" {
		return status, fmt.Errorf(message)
	}

	return status, nil

}

//GetSecondLegOrderID returns second LEg orders for a particular BO order.
func (ord Order) GetSecondLegOrderID(parentID string) (string, error) {
	orders, err := ord.KC.GetOrders()
	if err != nil {
		log.Println("Error retrieving parent Order ID")
		return "", err
	}

	for _, order := range orders {
		if order.ParentOrderID == parentID {
			return order.OrderID, nil
		}
	}

	return "", fmt.Errorf("No Second Leg Order found")

}

func getUnit32(str string) uint32 {
	// var a uint32
	u, _ := strconv.ParseUint(str, 10, 32)
	return uint32(u)
}
