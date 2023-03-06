package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendToDiscordWebhook(webhookURL string, message string) error {
	type discordMessage struct {
		Content string `json:"content"`
	}

	payload, err := json.Marshal(discordMessage{Content: message})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	return nil
}
