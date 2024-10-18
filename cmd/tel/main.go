package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/redsubmarine/tel/internal/command"
	"github.com/redsubmarine/tel/internal/config"
	"github.com/redsubmarine/tel/internal/telegram"
)

func setup() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter bot token: ")
	botToken, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Input error:", err)
		os.Exit(1)
	}
	botToken = strings.TrimSpace(botToken)

	// Retrieve chat ID
	chatID, err := telegram.GetChatID(botToken)
	if err != nil {
		fmt.Println("Failed to retrieve chat ID:", err)
		os.Exit(1)
	}

	// Save configuration
	cfg := &config.Config{
		BotToken: botToken,
		ChatID:   chatID,
	}

	if err := config.WriteConfig(cfg); err != nil {
		fmt.Println("Failed to save configuration:", err)
		os.Exit(1)
	}

	fmt.Println("Setup completed.")
}

func executeCommand(cfg *config.Config, args []string) {
	// Execute command
	output, exitStatus, err := command.ExecuteCommand(args)

	// Set status message
	var statusMsg string
	if exitStatus == 0 {
		statusMsg = "Command executed successfully."
	}
	if err != nil {
		statusMsg = fmt.Sprintf("Command failed with error: %v\n", err)
	}

	// Create Telegram message
	message := command.FormatMessage(args, statusMsg, output)

	// Send message asynchronously
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := telegram.SendMessage(cfg, message)
		if err != nil {
			fmt.Println("Failed to send Telegram message:", err)
		} else {
			fmt.Println("Telegram message sent successfully.")
		}
	}()
	wg.Wait()

	os.Exit(exitStatus)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tel <command> [arguments...] or tel setup")
		os.Exit(1)
	}

	commandArg := os.Args[1]

	if commandArg == "setup" {
		setup()
		os.Exit(0)
	}

	// Read configuration
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Unable to read the configuration file. Please run `tel setup` first.")
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Println("Invalid configuration:", err)
		os.Exit(1)
	}

	// Execute command and send notification
	executeCommand(cfg, os.Args[1:])
}
