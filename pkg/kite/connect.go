package kite

/*
Handling Kite connections and related operations here, like:
- Creating kite connection.
- Buy order
- Sell Order
*/

import (
	kiteconnect "github.com/zerodhatech/gokiteconnect"
)

const (
	apiKey    string = "c7b8qdb6dcwc9obc"
	apiSecret string = "ldxrh0w77x88zyhyuivbbr2svm2kml17"
)

//GetKiteAccessToken to retrieve access token for Kite API
func GetKiteAccessToken() {

}

// Connect creates a connection with the kite API.
func Connect() (*kiteconnect.Client, string) {

	// fmt.Println("Creating Saudgar!")
	// kiteSession := scrapper.NewWDSession()
	kc := kiteconnect.New(apiKey)
	// fmt.Println(kc.GetLoginURL())
	// authURL := kiteSession.GetKiteAuthURL(kc.GetLoginURL())
	// fmt.Println(authURL)
	// requestToken := "0w9iSbcyH57ofgk6nAfAdzeUNd7xf1LE"

	// //Get user details and access token
	// data, err := kc.GenerateSession(requestToken, apiSecret)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)

	// }

	// fmt.Println(data.AccessToken)
	// panic(1)
	accessToken := "6fbbf9hI4vM6T2fpYV1ylZEu7LtrFK0K"
	kc.SetAccessToken(accessToken)

	// holdings, err := kc.GetHoldings()
	// if err != nil {
	// 	fmt.Printf("Error getting margins: %v", err)
	// }
	// fmt.Println("holdings: ", holdings)
	// inst, _ := kc.GetInstruments()
	// fmt.Println(inst)
	// panic(1)

	return kc, accessToken

}
