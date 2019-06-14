package instrument

import "github.com/ashwanthkumar/slack-go-webhook"
import "fmt"

//SendAlerts sends alerts to slack.
func SendAlerts(message string) {
	webhookURL := "https://hooks.slack.com/services/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"

	// attachment1 := slack.Attachment{}
	// attachment1.AddField(slack.Field{Title: "Author", Value: "Ashwanth Kumar"}).AddField(slack.Field{Title: "Status", Value: "Completed"})
	// attachment1.AddAction(slack.Action{Type: "button", Text: "Book flights 🛫", Url: "https://flights.example.com/book/r123456", Style: "primary"})
	// attachment1.AddAction(slack.Action{Type: "button", Text: "Cancel", Url: "https://flights.example.com/abandon/r123456", Style: "danger"})
	payload := slack.Payload{
		Text:      message,
		Username:  "sp98",
		Channel:   "#kite-trade",
		IconEmoji: ":monkey_face:",
	}
	err := slack.Send(webhookURL, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}
