package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/phoenix-marie/core/internal/cli"
)

func main() {
	// Load environment
	godotenv.Load(".env.local")

	// Check for command mode
	if len(os.Args) > 1 {
		command := os.Args[1]
		args := strings.Join(os.Args[2:], " ")

		handler := cli.NewHandler()
		if err := handler.ExecuteCommand(command, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Interactive chat mode
	handler := cli.NewHandler()
	handler.StartInteractiveChat()
}

