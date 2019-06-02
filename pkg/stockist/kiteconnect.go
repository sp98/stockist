package stockist

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

// KiteConnect creates a connection with the kite API.
func KiteConnect() (*kiteconnect.Client, string) {

	// fmt.Println("Creating Saudgar!")
	// kiteSession := scrapper.NewWDSession()
	kc := kiteconnect.New(apiKey)
	// authURL := kiteSession.GetKiteAuthURL(kc.GetLoginURL())
	// fmt.Println(authURL)
	//requestToken := "ZDxeebXgqR9KoZYQ8xwdnzi6vRffR4wH"

	// //Get user details and access token
	// data, err := kc.GenerateSession(requestToken, apiSecret)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	return
	// }

	// fmt.Println(data.AccessToken)
	accessToken := "K3Aro6h07IXuQHZcXo6WbHCTsutiOY2V"
	kc.SetAccessToken(accessToken)

	// holdings, err := kc.GetHoldings()
	// if err != nil {
	// 	fmt.Printf("Error getting margins: %v", err)
	// }
	// fmt.Println("holdings: ", holdings)
	// inst, _ := kc.GetInstruments()
	// fmt.Println(inst)

	return kc, accessToken

}
