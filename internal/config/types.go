package config

// Config struct holds the contents of the configuration file.
type Config struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}
