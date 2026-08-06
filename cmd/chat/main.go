package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
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

	system := []anthropic.TextBlockParam{}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("User: ")

	for scanner.Scan() {
		prompt := scanner.Text()
		if prompt == "" {
			break
		}
		if lower := strings.ToLower(prompt); lower == "exit" || lower == "quit" || lower == "q" {
			break
		}
		messages = addUserMessage(messages, prompt)

		stream := chat(client, messages, model, system)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n--\nAssistant: ")
		response := ""
		for stream.Next() {
			data := stream.Current()
			if stream.Err() != nil {
				panic(stream.Err())
			}
			fmt.Printf("%s", data.Delta.Text)
			response = response + data.Delta.Text
		}

		messages = addAssistantMessage(messages, response)

		fmt.Printf("\n--\nUser: ")
	}
}

// chat sends the prompt to the API and returns the response.
func chat(
	client anthropic.Client,
	messages []anthropic.MessageParam,
	model string,
	system []anthropic.TextBlockParam,
) *ssestream.Stream[anthropic.MessageStreamEventUnion] {

	var temp param.Opt[float64]
	temp.Value = 1.0

	messageStream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens:   1024,
		Messages:    messages,
		Model:       model,
		System:      system,
		Temperature: temp,
	})

	return messageStream
}

// getResponse returns the response from the API to the console.
// Not needed for streaming, but useful for testing.
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
