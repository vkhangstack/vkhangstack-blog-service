package repository

import (
	"context"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

// CreateBotChatMessage inserts one bot<->user message into the chat log.
func (u *DB) CreateBotChatMessage(ctx context.Context, msg domain.BotChatMessage) (*domain.BotChatMessage, error) {
	msg.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&msg).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("bot chat message not saved: %v", err)
	}
	return &msg, nil
}

// CreateInboundBotChatMessage inserts an inbound message, deduped on the Zalo message_id via the
// partial unique index. A redelivered webhook (same message_id) hits ON CONFLICT DO NOTHING and
// affects zero rows, so the caller can skip reprocessing. Returns true only on a fresh insert.
func (u *DB) CreateInboundBotChatMessage(ctx context.Context, msg domain.BotChatMessage) (bool, error) {
	msg.ID = u.snowflakeNode.GenerateID()
	q := u.db.NewInsert().Model(&msg)
	if msg.MessageID != "" {
		q = q.On("CONFLICT (message_id) WHERE message_id <> '' DO NOTHING")
	}
	res, err := q.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("inbound bot chat message not saved: %v", err)
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// ListBotChatConversations returns one summary row per Zalo chat, newest activity first, with the
// most recent message and total message count. Total is the number of distinct conversations.
func (u *DB) ListBotChatConversations(ctx context.Context, limit, offset int) ([]*domain.BotChatConversation, int, error) {
	var total int
	if err := u.db.NewRaw(
		"SELECT COUNT(DISTINCT chat_id) FROM bot_chat_messages",
	).Scan(ctx, &total); err != nil {
		return nil, 0, fmt.Errorf("bot chat conversations count failed: %v", err)
	}

	var conversations []*domain.BotChatConversation
	err := u.db.NewRaw(
		`SELECT t.chat_id, t.user_id, t.user_name, t.last_text, t.last_sender, t.last_at, t.message_count
		 FROM (
		   SELECT DISTINCT ON (c.chat_id)
		     c.chat_id,
		     COALESCE(c.user_id, '') AS user_id,
		     c.text       AS last_text,
		     c.sender     AS last_sender,
		     c.created_at AS last_at,
		     (SELECT COUNT(*) FROM bot_chat_messages m WHERE m.chat_id = c.chat_id) AS message_count,
		     COALESCE((SELECT n.sender_name FROM bot_chat_messages n
		       WHERE n.chat_id = c.chat_id AND n.sender_name IS NOT NULL AND n.sender_name <> ''
		       ORDER BY n.created_at DESC LIMIT 1), '') AS user_name
		   FROM bot_chat_messages c
		   ORDER BY c.chat_id, c.created_at DESC
		 ) t
		 ORDER BY t.last_at DESC
		 LIMIT ? OFFSET ?`,
		limit, offset,
	).Scan(ctx, &conversations)
	if err != nil {
		return nil, 0, fmt.Errorf("bot chat conversations query failed: %v", err)
	}
	return conversations, total, nil
}

// ListBotChatMessages returns the messages of one conversation, oldest first (chat order).
func (u *DB) ListBotChatMessages(ctx context.Context, chatID string, limit, offset int) ([]*domain.BotChatMessage, int, error) {
	var messages []*domain.BotChatMessage
	total, err := u.db.NewSelect().
		Model(&messages).
		Where("bcm.chat_id = ?", chatID).
		Order("bcm.created_at ASC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("bot chat messages query failed: %v", err)
	}
	return messages, total, nil
}
