package orders

import "github.com/ashwanthkumar/slack-go-webhook"
import "fmt"

//SendAlerts sends alerts to slack.
func SendAlerts(message string) {
	webhookURL := "https://hooks.slack.com/services/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"
	payload := slack.Payload{
		Text:     message,
		Username: "sp98",
		Channel:  "#kite-trade",
	}
	err := slack.Send(webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}
