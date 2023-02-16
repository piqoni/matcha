package main

import (
	"context"
	"log"
	"time"

	readability "github.com/go-shiori/go-readability"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func getSummaryFromLink(url string) string {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to parse %s, %v\n", url, err)
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
	c := gogpt.NewClient(openaiApiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 60,
		Prompt:    text + " \n\nTl;dr",
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return ""
	}

	return resp.Choices[0].Text
}
