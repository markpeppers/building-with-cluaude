package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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

	messages := []anthropic.MessageParam{}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("User: ")

	for scanner.Scan() {
		prompt := scanner.Text()
		if prompt == "" {
			break
		}
		if lower := strings.ToLower(prompt); lower == "exit" || lower == "quit" {
			break
		}
		messages = addUserMessage(messages, prompt)

		responseText, err := chat(client, messages, model)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n--\nAssistant: %s\n", responseText)

		messages = addAssistantMessage(messages, responseText)

		fmt.Printf("\n--\nUser: ")
	}
}

// chat sends the prompt to the API and returns the response.
func chat(client anthropic.Client, messages []anthropic.MessageParam, model string) (string, error) {

	system := []anthropic.TextBlockParam{
		{
			Text: "You are Bugs Bunny.",
		},
	}

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages:  messages,
		Model:     model,
		System:    system,
	})
	if err != nil {
		return "", err
	}

	responseText := getResponse(message)
	return responseText, nil
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
