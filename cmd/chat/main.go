package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	key := os.Getenv("ANTHROPIC_API_KEY")

	client := anthropic.NewClient(
		option.WithAPIKey(key),
	)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(
					"What are quaternions in one sentence?",
				),
			),
		},
		Model: anthropic.ModelClaudeOpus5,
	})
	if err != nil {
		panic(err)
	}

	printResponse(message)
}

// printResponse prints the response from the API to the console.
func printResponse(response *anthropic.Message) {
	for _, block := range response.Content {
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok {
			fmt.Println(textBlock.Text)
		}
	}
}
