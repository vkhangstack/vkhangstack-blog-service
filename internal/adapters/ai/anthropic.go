package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

const defaultClaudeModel = "claude-opus-4-8"

// ClaudeGenerator generates quizzes via the Anthropic Messages API.
type ClaudeGenerator struct {
	client anthropic.Client
	model  anthropic.Model
}

func NewClaudeGenerator(apiKey, model string) *ClaudeGenerator {
	if model == "" {
		model = defaultClaudeModel
	}
	return &ClaudeGenerator{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  anthropic.Model(model),
	}
}

func (g *ClaudeGenerator) GenerateQuiz(ctx context.Context, req domain.GenerateQuizRequest) (*domain.CreateQuizRequest, error) {
	system, user := buildPrompt(req)

	adaptive := anthropic.ThinkingConfigAdaptiveParam{}
	resp, err := g.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     g.model,
		MaxTokens: 16000,
		Thinking:  anthropic.ThinkingConfigParamUnion{OfAdaptive: &adaptive},
		System: []anthropic.TextBlockParam{
			{Text: system},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude request failed: %v", err)
	}
	if resp.StopReason == anthropic.StopReasonRefusal {
		return nil, fmt.Errorf("claude declined the request")
	}

	var text string
	for _, block := range resp.Content {
		if variant, ok := block.AsAny().(anthropic.TextBlock); ok {
			text += variant.Text
		}
	}
	if text == "" {
		return nil, fmt.Errorf("claude returned no text content")
	}

	return parseQuizPayload(text, req)
}

func (g *ClaudeGenerator) GenerateFlashcards(ctx context.Context, req domain.GenerateFlashcardsRequest) (*domain.GenerateFlashcardsResult, error) {
	system, user := buildFlashcardPrompt(req)

	adaptive := anthropic.ThinkingConfigAdaptiveParam{}
	resp, err := g.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     g.model,
		MaxTokens: 16000,
		Thinking:  anthropic.ThinkingConfigParamUnion{OfAdaptive: &adaptive},
		System: []anthropic.TextBlockParam{
			{Text: system},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(user)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude request failed: %v", err)
	}
	if resp.StopReason == anthropic.StopReasonRefusal {
		return nil, fmt.Errorf("claude declined the request")
	}

	var text string
	for _, block := range resp.Content {
		if variant, ok := block.AsAny().(anthropic.TextBlock); ok {
			text += variant.Text
		}
	}
	if text == "" {
		return nil, fmt.Errorf("claude returned no text content")
	}

	return parseFlashcardPayload(text)
}
