package main

import (
	"context"
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
	model := os.Getenv("MODEL")
	if model == "" {
		model = "claude-opus-5"
	}

	client := anthropic.NewClient(
		option.WithAPIKey(key),
	)

	prompts := []string{
		"What are quaternions in one sentence?",
		"Write another sentence.",
	}

	messages := []anthropic.MessageParam{}

	messages = addUserMessage(messages, "You are Elmer Fudd, a character from Looney Tunes. You speak in a funny way and have a lisp. You are also very polite and kind.")

	for _, prompt := range prompts {
		messages = addUserMessage(messages, prompt)

		message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages:  messages,
			Model:     model,
		})
		if err != nil {
			panic(err)
		}

		responseText := getResponse(message)
		println(responseText)

		messages = addAssistantMessage(messages, responseText)
	}
}

// getResponse returns the response from the API to the console.
func getResponse(response *anthropic.Message) string {
	for _, block := range response.Content {
		if textBlock, ok := block.AsAny().(anthropic.TextBlock); ok {
			return textBlock.Text
		}
	}
	return ""
}

func addUserMessage(messages []anthropic.MessageParam, text string) []anthropic.MessageParam {
	messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(text)))
	return messages
}

func addAssistantMessage(messages []anthropic.MessageParam, text string) []anthropic.MessageParam {
	messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(text)))
	return messages
}
