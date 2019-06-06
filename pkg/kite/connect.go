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
	// panic(1)
	// authURL := kiteSession.GetKiteAuthURL(kc.GetLoginURL())
	// fmt.Println(authURL)
	// requestToken := "c7b8qdb6dcwc9obc&sess_id=ThmTF32Q"

	// //Get user details and access token
	// data, err := kc.GenerateSession(requestToken, apiSecret)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)

	// }

	// fmt.Println(data.AccessToken)
	// panic(1)
	accessToken := "b4xnjM305ErdwvQiFujg7kP199EQubeE"
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
