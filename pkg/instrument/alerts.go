package instrument

import "github.com/ashwanthkumar/slack-go-webhook"
import "fmt"

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
	if channel == openTrendChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == buyStockChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == sellStockChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == sensexTrendChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	} else if channel == shortSellStocksChannel {
		webhookURL = "https://hooks.slack.com/services/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"
	}

	return webhookURL

}
