package zalobot

import (
	"encoding/json"
	"errors"
	"fmt"

	zalobot "github.com/vkhangstack/go-zalo-bot"
	"github.com/vkhangstack/go-zalo-bot/types"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

type ZaloBotConfig struct {
	Token         string
	WebhookSecret string
	Endpoint      string
	DeepLink      string
}

type ZaloBotAdapter struct {
	bot       *zalobot.BotAPI
	botConfig ZaloBotConfig
}

func NewZaloBotAdapter(cfg ZaloBotConfig) (*ZaloBotAdapter, error) {
	bot, err := zalobot.New(cfg.Token, types.WithDebug(), types.WithEnvironment(types.Production))

	if err != nil {
		return nil, err
	}
	if cfg.WebhookSecret != "" {
		bot.SetWebhookSecretToken(cfg.WebhookSecret)
	}
	if cfg.Endpoint != "" {
		bot.DeleteWebhook()
		bot.SetWebhook(types.WebhookConfig{
			URL:         cfg.Endpoint,
			SecretToken: cfg.WebhookSecret,
			Certificate: "", // Optional: Provide a certificate if needed
		})
	}
	return &ZaloBotAdapter{bot: bot, botConfig: cfg}, nil
}

// SendMessage sends a text message to chatID (the Zalo user's extend_id).
func (z *ZaloBotAdapter) SendMessage(chatID, text string) error {
	_, err := z.bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   text,
	})
	return err
}

// zaloWebhookEnvelope mirrors the real Zalo Bot webhook payload shape, which wraps the
// update inside a "result" field (e.g. {"ok":true,"result":{"message":{...},"event_name":
// "message.text.received"}}). The SDK's own ParseUpdate only recognizes update fields at
// the top level or its own flat WebhookEvent shape, so it fails to parse this envelope and
// falls through to "unsupported event type" — we parse the real shape ourselves instead.
type zaloWebhookEnvelope struct {
	EventName string `json:"event_name"`
	Message   struct {
		Chat struct {
			ID       string `json:"id"`
			ChatType string `json:"chat_type"`
		} `json:"chat"`
		Date int64 `json:"date"`
		From struct {
			ID          string `json:"id"`
			DisplayName string `json:"display_name"`
			IsBot       bool   `json:"is_bot"`
		} `json:"from"`
		MessageID string `json:"message_id"`
		Text      string `json:"text"`
	} `json:"message"`
}

// ProcessWebhook validates and parses an incoming Zalo Bot webhook payload, returning the
// sender's Zalo user ID (extend_id) and the text they sent. Messages sent by the bot itself
// (from.is_bot == true) are rejected with ports.ErrZaloBotMessage so callers can ignore them.
func (z *ZaloBotAdapter) ProcessWebhook(payload []byte, signature string) (senderID, text string, err error) {
	if signature == "" || signature != z.botConfig.WebhookSecret {
		return "", "", fmt.Errorf("zalo webhook: invalid signature")
	}
	var envelope zaloWebhookEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return "", "", fmt.Errorf("zalo webhook: failed to parse payload: %w", err)
	}
	if envelope.Message.MessageID == "" || envelope.Message.Text == "" {
		return "", "", errors.New("zalo webhook payload has no message")
	}
	if envelope.Message.From.IsBot {
		return "", "", domain.ErrZaloBotMessage
	}
	return envelope.Message.From.ID, envelope.Message.Text, nil
}

func (z *ZaloBotAdapter) GetSecretToken() string {
	return z.botConfig.WebhookSecret
}

func (z *ZaloBotAdapter) GetFieldSecretToken() string {
	return z.bot.GetFieldSecretToken()
}

func (z *ZaloBotAdapter) GetDeepLink() string {
	return z.botConfig.DeepLink
}
