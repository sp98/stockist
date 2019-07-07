package notification

import "github.com/ashwanthkumar/slack-go-webhook"
import "fmt"

const (
	BuyStockChannel        = "#buy-stocks"
	SellStockChannel       = "#sell-stocks"
	OpenTrendChannel       = "#opening-trend"
	SensexTrendChannel     = "#sensex-trend"
	ShortSellStocksChannel = "#short-sell"
	TradeChannel           = "#trade"
	ErrorChannel           = "#errors"
	OpenLowHigh            = "#open-low-high"
)

//OpeningTrade attachment format
type OpeningTrade struct {
	Message       string
	Instrument    string
	Exchange      string
	Token         string
	Open          float64
	PreviousOpen  float64
	PreviousHigh  float64
	PreviousClose float64
	ColoreCode    string
}

//NewOT returns Opening Trade message
func (ot OpeningTrade) NewOT(message, instrument, exchange, color, token string, open, prevOpen, prevHigh, prevClose float64) OpeningTrade {
	otm := &OpeningTrade{
		Message:       message,
		Instrument:    instrument,
		Exchange:      exchange,
		Token:         token,
		ColoreCode:    color,
		Open:          open,
		PreviousOpen:  prevOpen,
		PreviousClose: prevClose,
		PreviousHigh:  prevHigh,
	}

	return *otm
}

//SendAlerts sends alerts to slack.
func SendAlerts(message, channel string) {
	webhookURL := getWebhook(channel)
	payload := slack.Payload{
		Text: message,
		// Channel: channel,
	}
	err := slack.Send(webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}

func getWebhook(channel string) string {
	var webhookURL string
	if channel == OpenTrendChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == BuyStockChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == SellStockChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == SensexTrendChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == ShortSellStocksChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == TradeChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == ErrorChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == OpenLowHigh {
		webhookURL = "https://hooks.slack.com/services/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"
	}

	return webhookURL

}
