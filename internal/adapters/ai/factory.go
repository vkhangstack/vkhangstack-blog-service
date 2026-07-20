package ai

import (
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

// NewQuizGenerator builds the QuizGenerator for the configured provider.
// Returns nil (no error) when no API key is configured — the feature is
// simply disabled in that case.
func NewQuizGenerator(cfg config.AIConfig) (ports.QuizGenerator, error) {
	if cfg.APIKey == "" {
		return nil, nil
	}

	switch cfg.Provider {
	case "claude", "anthropic", "":
		return NewClaudeGenerator(cfg.APIKey, cfg.Model), nil
	case "openai", "chatgpt":
		model := cfg.Model
		if model == "" {
			model = defaultOpenAIModel
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = defaultOpenAIBaseURL
		}
		return NewOpenAICompatGenerator(cfg.APIKey, model, baseURL), nil
	case "kimi", "moonshot":
		model := cfg.Model
		if model == "" {
			model = defaultKimiModel
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = defaultKimiBaseURL
		}
		return NewOpenAICompatGenerator(cfg.APIKey, model, baseURL), nil
	default:
		return nil, fmt.Errorf("unknown AI provider: %q (supported: claude, openai, kimi)", cfg.Provider)
	}
}

// NewFlashcardGenerator builds the FlashcardGenerator for the configured provider.
// Returns nil (no error) when no API key is configured — the feature is
// simply disabled in that case.
func NewFlashcardGenerator(cfg config.AIConfig) (ports.FlashcardGenerator, error) {
	if cfg.APIKey == "" {
		return nil, nil
	}

	switch cfg.Provider {
	case "claude", "anthropic", "":
		return NewClaudeGenerator(cfg.APIKey, cfg.Model), nil
	case "openai", "chatgpt":
		model := cfg.Model
		if model == "" {
			model = defaultOpenAIModel
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = defaultOpenAIBaseURL
		}
		return NewOpenAICompatGenerator(cfg.APIKey, model, baseURL), nil
	case "kimi", "moonshot":
		model := cfg.Model
		if model == "" {
			model = defaultKimiModel
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = defaultKimiBaseURL
		}
		return NewOpenAICompatGenerator(cfg.APIKey, model, baseURL), nil
	default:
		return nil, fmt.Errorf("unknown AI provider: %q (supported: claude, openai, kimi)", cfg.Provider)
	}
}
