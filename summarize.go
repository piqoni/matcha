package main

import (
	"context"
	"fmt"
	"time"

	readability "github.com/go-shiori/go-readability"
	openai "github.com/sashabaranov/go-openai"
)

func getSummaryFromLink(url string) string {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		fmt.Printf("Failed to parse %s, %v\n", url, err)
	}

	return summarize(article.TextContent)

}

func summarize(text string) string {
	// Not sending everything to preserve Openai tokens in case the article is too long
	maxCharactersToSummarize := 5000
	if len(text) > maxCharactersToSummarize {
		text = text[:maxCharactersToSummarize]
	}

	// Dont summarize if the article is too short
	if len(text) < 200 {
		return ""
	}
	clientConfig := openai.DefaultConfig(openaiApiKey)
	if openaiBaseURL != "" {
		clientConfig.BaseURL = openaiBaseURL
	}
	client := openai.NewClientWithConfig(clientConfig)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "Summarize the following text:",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}
