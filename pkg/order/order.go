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
	Target1       float64
	Target2       float64
	PreviousClose float64
	OpenPrice     float64
	Params        *kiteconnect.OrderParams
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

	for _, order := range BuyOrders {
		var sl float64
		var so float64
		if order.TransactionType == "BUY" {
			sl = order.Price - order.StopLoss
			if order.Target2 != 0 {
				so = order.Target2 - order.Price
			} else {
				so = order.Target1 - order.Price
			}

		} else if order.TransactionType == "SELL" {
			sl = order.StopLoss - order.Price
			if order.Target2 != 0 {
				so = order.Price - order.Target2
			} else {
				so = order.Price - order.Target1
			}
		}
		ord := &Order{
			KC:      kc,
			Token:   order.Token,
			Variety: order.Variety,
			Target1: order.Target1,
			Target2: order.Target2,
			Params: &kiteconnect.OrderParams{
				Product:         order.Product,
				Exchange:        order.Exchange,
				Tradingsymbol:   order.Symbol,
				OrderType:       order.OrderType,
				Price:           order.Price,
				TransactionType: order.TransactionType,
				Stoploss:        float64(sl),
				Squareoff:       float64(so),
				Quantity:        order.Quantity,
				Validity:        order.Validity,
			},
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
	for _, ord := range *orders {
		log.Printf("Order %+v", ord)
		log.Printf("Order Param %+v", ord.Params)
		log.Println("*************************")
	}

	c := make(chan string)

	for _, order := range *orders {
		go order.execute(c)
	}

	for i := 0; i < len(*orders); i++ {
		log.Println("Result: ", <-c)
	}

}

func (ord Order) execute(c chan string) {

	prevClose, err := ord.GetClosePrice()
	if err != nil {
		log.Printf("Error finding Previous day close Price for the Instrument : %+v", ord.Params.Tradingsymbol)
	}

	ord.PreviousClose = prevClose
	openPrice, err := ord.GetOpenPrice()

	if err != nil {
		log.Printf("Error finding Opening Price for the Instrument : %+v", ord.Params.Tradingsymbol)
	}

	ord.OpenPrice = openPrice

	orderResp, err := ord.placeOrder()
	if err != nil {
		msg := fmt.Sprintf("Error Placing Order: %+v", err)
		log.Println(msg)
		alerts.SendAlerts(msg, alerts.ErrorChannel)
		return
	}

	status, err := ord.getOrderStatus(orderResp.OrderID)

	if status == "REJECTED" {
		//Send Notification
		ord.notfiyOrderRejection(orderResp.OrderID)
		c <- "REJECTED"
		return
	} else if status == "COMPLETE" {
		ord.notfiyOrderCompletion(orderResp.OrderID)
	}

	err = ord.exitOrder(orderResp.OrderID)
	if err != nil {
		msg := fmt.Sprintf("Order Exit Failed")
		alerts.SendAlerts(msg, alerts.ErrorChannel)
	}

	c <- "complete"
}

func (ord Order) placeOrder() (*kiteconnect.OrderResponse, error) {

	for {
		time.Sleep(500 * time.Millisecond)
		ltp, err := ord.GetLastTradingPrice()
		if err != nil {
			log.Println("Error finding LTP for : ", ord.Params.Tradingsymbol)
			//return orderIDs, err
		}
		//log.Println("LTP Price : ", ltp)
		if ord.Params.TransactionType == "BUY" { //ord.OpenPrice < ord.PreviousClose &&
			if ord.OpenPrice < ord.Params.Price && ltp >= ord.Params.Price {
				orderResp, err := ord.ExecuteOrder()
				if err != nil {
					return nil, err
				}
				return orderResp, nil
			}

		} else if ord.Params.TransactionType == "SELL" { // ord.OpenPrice > ord.PreviousClose
			if ord.OpenPrice > ord.Params.Price && ltp <= ord.Params.Price {
				orderResp, err := ord.ExecuteOrder()
				if err != nil {
					return nil, err
				}
				return orderResp, nil
			}
		}

		log.Printf("Waiting for %s order to be placed on %s", ord.Params.TransactionType, ord.Params.Tradingsymbol)

	}

	//return nil, fmt.Errorf("Order %+v didn't get execute", ord)

}

func (ord Order) exitOrder(orderID string) error {
	//Exit Order when:
	// 1 - Profit goes less than 400
	// 2 - LTP greater than Target 2
	// If ltp is greate than Target 1 then update order the order with stop lose of Target 1
	log.Printf("Exiting Order  : %+v", orderID)
	//ord.exit(orderIDs)

	for {
		time.Sleep(500 * time.Millisecond)
		urp, _ := ord.GetUnRealisedProfit()
		if urp < -5 { //400
			err := ord.exit(orderID)
			if err != nil {
				return err
			}
		}
		ltp, _ := ord.GetLastTradingPrice()
		if ltp > ord.Target1 {
			msg := fmt.Sprintf("Target1 Achieved: \n Instrument:%s \nOrder: %s", ord.Params.Tradingsymbol, ord.Params.TransactionType)
			alerts.SendAlerts(msg, alerts.TradeChannel)
			//update order to change trigger price
		}

		if ltp > ord.Target2 {
			err := ord.exit(orderID)
			if err != nil {
				return err
			}
		}

		log.Printf("Waiting for %s order to be placed on %s", ord.Params.TransactionType, ord.Params.Tradingsymbol)
	}

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
		log.Printf("Error exeting the order - %+v", err)
		return err
	}

	return nil
}

//Notify events on Slack
func (ord Order) notfiyOrderRejection(id string) {

	msg := fmt.Sprintf("ORDER REJECTED \nStock: %s \nExchange: %s \nOrderID: %s \n%s", ord.Params.Tradingsymbol, ord.Params.Exchange, id, separation)
	alerts.SendAlerts(msg, alerts.ErrorChannel)

}

//Notify events on Slack
func (ord Order) notfiyOrderCompletion(id string) {
	msg := fmt.Sprintf("ORDER COMPLETED \nStock: %s \nExchange: %s \nOrderID: %s \n%s", ord.Params.Tradingsymbol, ord.Params.Exchange, id, separation)
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
