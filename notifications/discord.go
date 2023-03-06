package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type discordMessage struct {
	Content string `json:"content"`
}

func SendToDiscordWebhook(webhookURL string, message string) error {
	payload, err := json.Marshal(discordMessage{Content: message})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	return nil
}
