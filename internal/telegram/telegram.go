package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/redsubmarine/tel/internal/config"
)

// SendMessage sends a message to Telegram.
func SendMessage(cfg *config.Config, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.BotToken)
	payload := fmt.Sprintf("chat_id=%s&text=%s&parse_mode=Markdown", cfg.ChatID, message)
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message: %s", string(body))
	}

	return nil
}

// GetChatID waits for the user to send a `start` message to the bot and returns the chat ID.
func GetChatID(botToken string) (string, error) {
	var chatID string
	offset := 0

	fmt.Println("Please send a `start` message to the bot...")

	for {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=100&offset=%d", botToken, offset)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching updates:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error reading response:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var result struct {
			Ok     bool     `json:"ok"`
			Result []Update `json:"result"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("JSON parsing error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range result.Result {
			if strings.ToLower(update.Message.Text) == "start" {
				chatID = fmt.Sprintf("%d", update.Message.Chat.ID)
				fmt.Printf("Chat ID found: %s\n", chatID)
				return chatID, nil
			}
			offset = update.UpdateID + 1
		}

		time.Sleep(1 * time.Second)
	}
}
