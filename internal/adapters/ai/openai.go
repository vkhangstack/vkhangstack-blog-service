package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

const (
	defaultOpenAIBaseURL = "https://api.openai.com/v1"
	defaultOpenAIModel   = "gpt-4o"
	defaultKimiBaseURL   = "https://api.moonshot.ai/v1"
	defaultKimiModel     = "kimi-k2-0711-preview"
)

// OpenAICompatGenerator generates quizzes via any OpenAI-compatible
// chat-completions API — used for both ChatGPT (OpenAI) and Kimi (Moonshot).
type OpenAICompatGenerator struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

func NewOpenAICompatGenerator(apiKey, model, baseURL string) *OpenAICompatGenerator {
	return &OpenAICompatGenerator{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (g *OpenAICompatGenerator) GenerateQuiz(ctx context.Context, req domain.GenerateQuizRequest) (*domain.CreateQuizRequest, error) {
	system, user := buildPrompt(req)

	body, err := json.Marshal(chatRequest{
		Model: g.model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ai request failed: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ai response read failed: %v", err)
	}

	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("ai response is not valid JSON: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg := string(raw)
		if parsed.Error != nil {
			msg = parsed.Error.Message
		}
		return nil, fmt.Errorf("ai request failed (%d): %s", resp.StatusCode, msg)
	}
	if len(parsed.Choices) == 0 || parsed.Choices[0].Message.Content == "" {
		return nil, fmt.Errorf("ai returned no content")
	}

	return parseQuizPayload(parsed.Choices[0].Message.Content, req)
}
