package order

import (
	"fmt"
	"log"
	"time"

	kiteconnect "github.com/sp98/gokiteconnect"
	alerts "github.com/stockist/pkg/notification"
)

const (
	separation = "*********************************************************"
)

//Order struct
type Order struct {
	KC            *kiteconnect.Client //Kite Trading Client
	Variety       string
	Token         string //Instruement token
	Symbol        string
	Exchange      string
	PreviousClose float64
	OpenPrice     float64
	Status        string
	OrderID       string
	BuyParams     *BuyParams
	SellParams    *SellParams
}

//BuyParams are parameters when buying a stock
type BuyParams struct {
	Target1 float64
	Target2 float64
	Params  *kiteconnect.OrderParams
}

//SellParams are parameters when Selling a stock
type SellParams struct {
	Target1 float64
	Target2 float64
	Params  *kiteconnect.OrderParams
}

//Status struct to tract the status of an order.
type Status struct {
	Symbol     string
	Exchange   string
	Status     string
	profitLoss float64
}

//RejectedOrder holds all the rejected orders and message
type RejectedOrder struct {
	OrderID string
	Message string
}

//New returns the list of order to be executed
func New() *[]Order {
	orderList := []Order{}
	kc := getConnection()
	bParams := &BuyParams{}
	sParams := &SellParams{}

	for _, order := range BuyOrders {
		var bsl float64
		var bso float64
		var ssl float64
		var sso float64

		bsl = order.BuyPrice - order.BuyStopLoss
		if order.BuyTarget2 != 0 {
			bso = order.BuyTarget2 - order.BuyPrice
		} else {
			bso = order.BuyTarget1 - order.BuyPrice
		}

		ssl = order.SellStopLoss - order.SellPrice
		if order.SellTarget2 != 0 {
			sso = order.SellPrice - order.SellTarget2
		} else {
			sso = order.SellPrice - order.SellTarget1
		}

		if order.BuyPrice == 0 {
			bParams = nil
		} else {
			bParams =
				&BuyParams{
					Target1: order.BuyTarget1,
					Target2: order.BuyTarget2,
					Params: &kiteconnect.OrderParams{
						Product:         order.Product,
						Exchange:        order.Exchange,
						Tradingsymbol:   order.Symbol,
						OrderType:       order.OrderType,
						Price:           order.BuyPrice,
						TransactionType: "BUY",
						Stoploss:        float64(bsl),
						Squareoff:       float64(bso),
						Quantity:        order.Quantity,
						Validity:        order.Validity,
					},
				}

		}

		if order.SellPrice == 0 {
			sParams = nil
		} else {
			sParams = &SellParams{
				Target1: order.SellTarget1,
				Target2: order.SellTarget2,
				Params: &kiteconnect.OrderParams{
					Product:         order.Product,
					Exchange:        order.Exchange,
					Tradingsymbol:   order.Symbol,
					OrderType:       order.OrderType,
					Price:           order.SellPrice,
					TransactionType: "SELL",
					Stoploss:        float64(ssl),
					Squareoff:       float64(sso),
					Quantity:        order.Quantity,
					Validity:        order.Validity,
				},
			}

		}
		ord := &Order{
			KC:         kc,
			Token:      order.Token,
			Symbol:     order.Symbol,
			Exchange:   order.Exchange,
			Variety:    order.Variety,
			BuyParams:  bParams,
			SellParams: sParams,
		}
		orderList = append(orderList, *ord)

	}

	return &orderList
}

//Start starts order execution
func Start() {
	orders := New()
	if len(*orders) == 0 {
		msg := "Error: No orders to execute as order list is empty!"
		log.Println(msg)
		alerts.SendAlerts(msg, alerts.ErrorChannel)
		return
	}

	log.Printf("Orders to be executed today: %+v", len(*orders))
	// for _, ord := range *orders {
	// 	log.Printf("Order %+v", ord)
	// 	log.Printf("Order Buy Param %+v", ord.BuyParams)
	// 	log.Printf("Order Buy Order Params %+v", ord.BuyParams.Params)
	// 	log.Printf("Order Sell Param %+v", ord.SellParams)
	// 	log.Printf("Order Sell Order Param %+v", ord.SellParams.Params)
	// 	log.Println("*************************")
	// }

	c := make(chan string)

	for i, order := range *orders {
		ord := &order

		//Verify if the order is already created
		posCreated := ord.PositionCreated()

		if posCreated {
			//Check if there are pending BO second leg orders for the stock
			pendingOrders, tradeType, parentOrderID := ord.PendingBOOrders()
			ord.OrderID = parentOrderID
			if pendingOrders {
				log.Printf("Position already created but pending BO orders are available for Stock - %s", ord.Symbol)
				if tradeType == "SELL" {
					ord.Status = "BOUGHT"
				} else if tradeType == "BUY" {
					ord.Status = "SOLD"
				}

			} else {
				log.Printf("Position already created and no pending BO orders are available for Stock - %s", ord.Symbol)
				if i == len(*orders)-1 {
					return //becuase no more orders are there to execute
				}
				continue
			}
		}

		go ord.execute(c)

	}

	for i := 0; i < len(*orders); i++ {
		log.Println("Result: ", <-c)
	}

}

func (ord Order) execute(c chan string) {
	var orderID string
	prevClose, err := ord.GetClosePrice()
	if err != nil {
		log.Printf("Error finding Previous day close Price for the Instrument : %+v", ord.Symbol)
	}

	ord.PreviousClose = prevClose
	openPrice, err := ord.GetOpenPrice()

	if err != nil {
		log.Printf("Error finding Opening Price for the Instrument : %+v", ord.Symbol)
	}

	ord.OpenPrice = openPrice

	//If Position is already created for this Stock, then skip placing new order.
	if ord.Status == "BOUGHT" || ord.Status == "SOLD" {
		orderID = ord.OrderID
	} else {
		orderResp, err := ord.placeOrder()
		if err != nil {
			msg := fmt.Sprintf("Error Placing Order: %+v", err)
			log.Println(msg)
			alerts.SendAlerts(msg, alerts.ErrorChannel)
			c <- "complete"
			return
		}
		orderID = orderResp.OrderID
	}

	if len(orderID) == 0 {
		msg := fmt.Sprintf("Failure: Order ID not received for the stock - %s. EXIT THE BO MANUALLY", ord.Symbol)
		alerts.SendAlerts(msg, alerts.ErrorChannel)
		c <- "FAILED"
		return
	}

	status, err := ord.getOrderStatus(orderID)
	log.Println("Order status -- ", status)
	if status == "REJECTED" {
		log.Println("Order rejected for: ", ord.Symbol)
		ord.notfiyOrderRejection(orderID)
		c <- "REJECTED"
		return
	} else if status == "COMPLETE" {
		ord.notfiyOrderCompletion(orderID)
	}

	err = ord.exitOrder(orderID)
	if err != nil {
		msg := fmt.Sprintf("Order Exit Failed")
		alerts.SendAlerts(msg, alerts.ErrorChannel)
		c <- "FAILED"
	}

	c <- "complete"
}

func (ord *Order) placeOrder() (*kiteconnect.OrderResponse, error) {

	//If order.status is not empty then get the existing orders Id from Position

	for {

		time.Sleep(500 * time.Millisecond)
		ltp, err := ord.GetLastTradingPrice()
		if err != nil {
			log.Println("Error finding LTP for : ", ord.Symbol)
			//return orderIDs, err
		}

		//log.Println("LTP Price : ", ltp)
		if ord.BuyParams != nil {
			if ord.OpenPrice < ord.BuyParams.Params.Price && ltp >= ord.BuyParams.Params.Price {
				orderResp, err := ord.ExecuteOrder("BUY")
				ord.Status = "BOUGHT"
				if err != nil {
					return nil, err
				}

				return orderResp, nil
			}
		}

		if ord.SellParams != nil {
			if ord.OpenPrice > ord.SellParams.Params.Price && ltp <= ord.SellParams.Params.Price {
				orderResp, err := ord.ExecuteOrder("SELL")
				if err != nil {
					return nil, err
				}
				ord.Status = "SOLD"
				return orderResp, nil
			}

		}

		log.Printf("Waiting for order to be placed on %s", ord.Symbol)

	}

	//return nil, fmt.Errorf("Order %+v didn't get execute", ord)

}

func (ord *Order) exitOrder(orderID string) error {
	//Exit Order when:
	// 1 - Profit goes less than 400
	// 2 - LTP greater than Target 2
	// If ltp is greate than Target 1 then update order the order with stop lose of Target 1
	log.Printf("Exiting Order  : %+v", orderID)
	//ord.exit(orderIDs)

	for {
		time.Sleep(500 * time.Millisecond)

		if ord.Status == "BOUGHT" {
			ord.exitBuyOrder(orderID)
		} else if ord.Status == "SOLD" {
			ord.exitSellOrder(orderID)
		}

		log.Printf("Waiting for %s order to be placed on %s", "SELL", ord.Symbol)
	}

}

func (ord Order) exitBuyOrder(orderID string) error {
	urp, _ := ord.GetUnRealisedProfit()

	if urp < -700 { //400
		err := ord.exit(orderID)
		if err != nil {
			return err
		}
	}
	ltp, _ := ord.GetLastTradingPrice()
	if ltp >= ord.BuyParams.Target1 {
		msg := fmt.Sprintf("Target1 Achieved: \n Instrument:%s ", ord.Symbol)
		alerts.SendAlerts(msg, alerts.TradeChannel)
		//update order to change trigger price
	}

	if ltp >= ord.BuyParams.Target2 {
		err := ord.exit(orderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ord Order) exitSellOrder(orderID string) error {
	urp, _ := ord.GetUnRealisedProfit()

	if urp < -700 { //400
		err := ord.exit(orderID)
		if err != nil {
			return err
		}
	}
	ltp, _ := ord.GetLastTradingPrice()
	if ltp <= ord.SellParams.Target1 {
		msg := fmt.Sprintf("Target1 Achieved: \n Instrument:%s", ord.Symbol)
		alerts.SendAlerts(msg, alerts.TradeChannel)
		//update order to change trigger price
	}

	if ltp <= ord.SellParams.Target2 {
		err := ord.exit(orderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ord Order) exit(parentOrderID string) error {
	secondLegOrderID, err := ord.GetSecondLegOrderID(parentOrderID)
	if err != nil {
		log.Printf("Error getting second Leg Order Id - %+v", err)
		return err
	}
	log.Printf("Parent Order ID %+v", parentOrderID)
	log.Printf("Second Leg Order ID %+v", secondLegOrderID)
	_, err = ord.KC.ExitOrder(ord.Variety, secondLegOrderID, &parentOrderID)
	if err != nil {
		log.Printf("Error executing the order - %+v", err)
		return err
	}

	return nil
}

//Notify events on Slack
func (ord Order) notfiyOrderRejection(id string) {

	msg := fmt.Sprintf("ORDER REJECTED \nStock: %s \nExchange: %s \nOrderID: %s \n%s", ord.Symbol, ord.Exchange, id, separation)
	alerts.SendAlerts(msg, alerts.ErrorChannel)

}

//Notify events on Slack
func (ord Order) notfiyOrderCompletion(id string) {
	msg := fmt.Sprintf("ORDER COMPLETED \nStock: %s \nExchange: %s \nOrderID: %s \n%s", ord.Symbol, ord.Exchange, id, separation)
	alerts.SendAlerts(msg, alerts.TradeChannel)

}

/*
Buy Call:
1. Open should be lower than previous close
2. LTP has reached the Buy Price.




Sell Call:
1. Open should be higher than Previous Close
2. LTP has reached the Buy Price
*/
