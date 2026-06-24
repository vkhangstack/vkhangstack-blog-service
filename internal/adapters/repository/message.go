package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func (m *DB) CreateMessage(userID string, message domain.Message) error {
	ctx := context.Background()
	message = domain.Message{
		ID:     m.snowflakeNode.GenerateID(),
		UserID: userID,
		Body:   message.Body,
	}
	_, err := m.db.NewInsert().Model(&message).Exec(ctx)
	if err != nil {
		return fmt.Errorf("messages not saved: %v", err)
	}
	return nil
}

func (m *DB) ReadMessage(id string) (*domain.Message, error) {
	ctx := context.Background()
	message := &domain.Message{}
	err := m.db.NewSelect().Model(message).Where("id = ?", id).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("message not found")
	}
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *DB) ReadMessages() ([]*domain.Message, error) {
	ctx := context.Background()
	var messages []*domain.Message
	err := m.db.NewSelect().Model(&messages).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("messages not found: %v", err)
	}
	return messages, nil
}

func (m *DB) UpdateMessage(id string, message domain.Message) error {
	ctx := context.Background()
	res, err := m.db.NewUpdate().Model(&message).OmitZero().Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("message not found")
	}
	return nil
}

func (m *DB) DeleteMessage(id string) error {
	ctx := context.Background()
	message := &domain.Message{}
	res, err := m.db.NewDelete().Model(message).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("message not found")
	}
	return nil
}
