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
	// requestToken := "lP31tDVdLIrtS8NL6avpV8r8U2d1DHwz"

	// //Get user details and access token
	// data, err := kc.GenerateSession(requestToken, apiSecret)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)

	// }

	// fmt.Println(data.AccessToken)
	// panic(1)
	accessToken := "U5fHzFXC7JHS4T5Ya9XfB2lResI5qrBv"
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
