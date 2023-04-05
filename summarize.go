package main

import (
	"context"
	"fmt"
	"strings"
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
	c := openai.NewClient(openaiApiKey)

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 60,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text + " \n\nTl;dr",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	// append ... if text does not end with .
	if !strings.HasSuffix(resp.Choices[0].Message.Content, ".") {
		resp.Choices[0].Message.Content += "..."
	}

	return resp.Choices[0].Message.Content
}
