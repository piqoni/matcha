package main

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type LLMClient struct {
	client *openai.Client
	config *Config
}

func NewLLMClient(cfg *Config) *LLMClient {
	cConfig := openai.DefaultConfig(cfg.OpenAIKey)
	if cfg.OpenAIBaseURL != "" {
		cConfig.BaseURL = cfg.OpenAIBaseURL
	}
	return &LLMClient{
		client: openai.NewClientWithConfig(cConfig),
		config: cfg,
	}
}

func (l *LLMClient) Summarize(text string) string {
	fmt.Println("summarize invoked")
	if l == nil || l.client == nil {
		return ""
	}

	const maxCharactersToSummarize = 5000
	const minCharactersToSummarize = 200

	if len(text) > maxCharactersToSummarize {
		text = text[:maxCharactersToSummarize]
	}

	// Don't summarize if the article is too short
	if len(text) < minCharactersToSummarize {
		return ""
	}

	prompt := l.config.SummaryPrompt
	if prompt == "" {
		prompt = "Summarize the following text:"
	}

	model := l.config.OpenAIModel
	if model == "" {
		model = openai.GPT3Dot5Turbo // TODO: change this
	}

	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: text,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Summarization error: %v\n", err)
		return ""
	}

	if len(resp.Choices) == 0 {
		return ""
	}

	return resp.Choices[0].Message.Content
}

func (l *LLMClient) Analyze(articles []string) string {
	if l == nil || len(articles) == 0 {
		return ""
	}

	model := l.config.AnalystModel
	if model == "" {
		model = openai.GPT4o
	}

	prompt := fmt.Sprintf("%s\n\n%s", l.config.AnalystPrompt, strings.Join(articles, "\n"))

	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
		},
	)
	if err != nil {
		fmt.Printf("Analysis failed: %v\n", err)
		return ""
	}
	return resp.Choices[0].Message.Content
}
