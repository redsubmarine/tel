package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// ExecuteCommand runs the given command, displaying real-time output to the terminal
// and capturing the complete output to return it.
func ExecuteCommand(args []string) (string, int, error) {
	if len(args) == 0 {
		return "", 1, fmt.Errorf("no command provided")
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(cmdName, cmdArgs...)

	// Connect stdout and stderr to pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", 1, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", 1, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", 1, fmt.Errorf("failed to start command: %v", err)
	}

	// Buffer to capture the output
	var outputBuf bytes.Buffer

	// Function to process the output from a pipe
	processPipe := func(pipe io.ReadCloser) {
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)                  // Real-time output
			outputBuf.WriteString(line + "\n") // Save to buffer
		}
	}

	// Process both stdout and stderr simultaneously
	go processPipe(stdoutPipe)
	go processPipe(stderrPipe)

	// Wait for the command to complete
	err = cmd.Wait()

	exitStatus := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitStatus = exitError.ExitCode()
		} else {
			exitStatus = 1
		}
	}

	return outputBuf.String(), exitStatus, err
}

// FormatMessage formats the message to be sent to Telegram.
func FormatMessage(args []string, statusMsg, output string) string {
	cmdName := args[0]
	cmdArgs := args[1:]
	cmdStr := fmt.Sprintf("%s %s", cmdName, strings.Join(cmdArgs, " "))

	// If the output is too large, consider including a summary or partial output.
	// Here, the full output is included.
	message := fmt.Sprintf("Execution completed.\n`%s`\nStatus: %s\n\nOutput:\n```\n%s\n```",
		cmdStr,
		statusMsg,
		output,
	)
	return message
}
