package services

import (
	"context"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// ChatService owns the bot<->user message log that backs the backoffice Chats inbox: it records
// every inbound/outbound message and lets an admin push a manual reply to a user through the bot.
type ChatService struct {
	repo        ports.ChatRepository
	settingRepo ports.NotificationSettingRepository
	bot         ports.ZaloBotClient
}

func NewChatService(
	repo ports.ChatRepository,
	settingRepo ports.NotificationSettingRepository,
	bot ports.ZaloBotClient,
) *ChatService {
	return &ChatService{repo: repo, settingRepo: settingRepo, bot: bot}
}

// resolveUserID best-effort maps a Zalo chat ID to the linked internal account; empty when the
// Zalo channel is not linked. Errors are swallowed so message logging never blocks on it.
func (s *ChatService) resolveUserID(ctx context.Context, chatID string) string {
	setting, err := s.settingRepo.GetNotificationSettingByChannelToken(ctx, domain.NotificationChannelZaloBot, chatID)
	if err != nil || setting == nil {
		return ""
	}
	return setting.UserID
}

// RecordInbound stores a message the user sent to the bot (with the Zalo sender's display name),
// deduped on the Zalo messageID so a redelivered webhook is a no-op. Returns true only when a new
// message was recorded, so the caller processes each Zalo message exactly once.
func (s *ChatService) RecordInbound(ctx context.Context, chatID, messageID, senderName, text string) bool {
	if chatID == "" || text == "" {
		return false
	}
	inserted, err := s.repo.CreateInboundBotChatMessage(ctx, domain.BotChatMessage{
		ChatID:     chatID,
		MessageID:  messageID,
		UserID:     s.resolveUserID(ctx, chatID),
		SenderName: senderName,
		Direction:  domain.BotChatDirectionInbound,
		Sender:     domain.BotChatSenderUser,
		Text:       text,
	})
	if err != nil {
		logger.Log.WithError(err).WithField("chat_id", chatID).Error("ChatService: failed to record inbound message")
		return false
	}
	return inserted
}

// RecordOutbound stores a message the bot/admin sent to the user. sender is BotChatSenderBot for
// automated command replies or BotChatSenderAdmin for manual console replies.
func (s *ChatService) RecordOutbound(ctx context.Context, chatID, sender, text string) {
	s.record(ctx, chatID, "", domain.BotChatDirectionOutbound, sender, text)
}

func (s *ChatService) record(ctx context.Context, chatID, senderName, direction, sender, text string) {
	if chatID == "" || text == "" {
		return
	}
	_, err := s.repo.CreateBotChatMessage(ctx, domain.BotChatMessage{
		ChatID:     chatID,
		UserID:     s.resolveUserID(ctx, chatID),
		SenderName: senderName,
		Direction:  direction,
		Sender:     sender,
		Text:       text,
	})
	if err != nil {
		logger.Log.WithError(err).WithField("chat_id", chatID).Error("ChatService: failed to record message")
	}
}

// SendReply delivers a manual admin reply to the user through the bot and records it. Returns an
// error only when the bot delivery fails; the caller surfaces that to the backoffice.
func (s *ChatService) SendReply(ctx context.Context, chatID, text string) error {
	if err := s.bot.SendMessage(chatID, text); err != nil {
		logger.Log.WithError(err).WithField("chat_id", chatID).Error("ChatService: failed to send admin reply")
		return err
	}
	s.RecordOutbound(ctx, chatID, domain.BotChatSenderAdmin, text)
	return nil
}

// ListConversations returns the paginated Chats inbox (one row per Zalo chat).
func (s *ChatService) ListConversations(ctx context.Context, page, limit int) ([]*domain.BotChatConversation, int, error) {
	limit, offset := normalizePage(page, limit)
	return s.repo.ListBotChatConversations(ctx, limit, offset)
}

// ListMessages returns the paginated messages of a single conversation, oldest first.
func (s *ChatService) ListMessages(ctx context.Context, chatID string, page, limit int) ([]*domain.BotChatMessage, int, error) {
	limit, offset := normalizePage(page, limit)
	return s.repo.ListBotChatMessages(ctx, chatID, limit, offset)
}

func normalizePage(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return limit, (page - 1) * limit
}
