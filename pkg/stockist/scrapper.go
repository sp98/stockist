package stockist

/*
Handle all the scrapping activities in this file like:
1. Get access token for Kite Trading API

*/
import (
	"fmt"
	"os"
	"strings"
	"time"

	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
)

const (
	seleniumPath = "/drivers/selenium-server-standalone-3.14.0.jar"
	//geckoDriverPath = "drivers/geckodriver-v0.23.0-linux64"
	geckoDriverPath = "/drivers/chromedriver-linux64-2.42"
	port            = 8080
)

const (
	loginUserID = "//input[@type='text']"
	loginPswd   = "//div[@class='login-form']//following::input[@type='password']"
	loginBtn    = "//button[@type='submit']"
	loginPin    = "//label[text()='PIN']/following::input[@type='password']"
	continueBtn = "//button[contains(text(), 'Continue')]"
	authoPage   = "//span[text()='This site can’t be reached']"
)

//Scrapper is the struct for a Webdriver session.
type Scrapper struct {
	Service *selenium.Service
	Wd      selenium.WebDriver
}

//NewWDSession creates a new Webdriver session
func NewWDSession() Scrapper {

	path, _ := os.Getwd()
	seleniumPath := path + seleniumPath
	geckoDriverPath := path + geckoDriverPath

	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	wd.SetAsyncScriptTimeout(10 * time.Second)
	wd.SetImplicitWaitTimeout(10 * time.Second)

	scrapper := Scrapper{Service: service, Wd: wd}
	return scrapper
}

//GetKiteAuthURL retrieves Kit auth URL
func (scrap Scrapper) GetKiteAuthURL(url string) string {

	// defer service.Stop()
	// // Connect to the WebDriver instance running locally.
	// caps := selenium.Capabilities{"browserName": "firefox"}
	// wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	// if err != nil {
	// 	panic(err)
	// }
	// defer wd.Quit()
	// wd.SetAsyncScriptTimeout(10 * time.Second)
	// wd.SetImplicitWaitTimeout(10 * time.Second)

	if err := scrap.Wd.Get(url); err != nil {
		panic(err)
	}

	// Find User name field
	userID, err := scrap.Wd.FindElement(selenium.ByXPATH, loginUserID)
	if err != nil {
		panic(err)
	}
	// Enter User name
	if err := userID.SendKeys("PJ7746"); err != nil {
		panic(err)
	}

	// Find password field.
	pswd, err := scrap.Wd.FindElement(selenium.ByXPATH, loginPswd)
	if err != nil {
		panic(err)
	}
	// Enter Password
	if err := pswd.SendKeys("$Ankump1l"); err != nil {
		panic(err)
	}

	// Find the login button
	login, err := scrap.Wd.FindElement(selenium.ByXPATH, loginBtn)
	if err != nil {
		panic(err)
	}
	// Click on the login button.
	if err := login.Click(); err != nil {
		panic(err)
	}

	pin, err := scrap.Wd.FindElement(selenium.ByXPATH, loginPin)
	if err != nil {
		println("I have reached here!")
		panic(err)
	}
	println("Pin Found!")
	// Enter Password
	if err := pin.SendKeys("109800"); err != nil {
		panic(err)
	}

	continbtn, err := scrap.Wd.FindElement(selenium.ByXPATH, continueBtn)
	if err != nil {
		panic(err)
	}
	// Click on the login button.
	if err := continbtn.Click(); err != nil {
		panic(err)
	}

	if err := scrap.Wd.WaitWithTimeoutAndInterval(conditions.ElementIsLocatedAndVisible(selenium.ByXPATH, authoPage), 30*time.Second, 5*time.Second); err != nil {
		println("pin not found ")
		panic(err)
	}

	authURL, err := scrap.Wd.CurrentURL()
	if err != nil {
		panic(err)
	}

	return authURL

}

//ParseAuthURL fetchs access token from the auth url
func ParseAuthURL(url string) string {
	accessToken := strings.TrimLeft(strings.TrimRight(url, "&action"), "?request_token")
	return accessToken
}
