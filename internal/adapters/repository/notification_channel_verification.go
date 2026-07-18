package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func (u *DB) CreateChannelVerification(ctx context.Context, verification domain.NotificationChannelVerification) (*domain.NotificationChannelVerification, error) {
	verification.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&verification).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("channel verification not saved: %v", err)
	}
	return &verification, nil
}

func (u *DB) GetPendingChannelVerification(ctx context.Context, channel domain.NotificationChannelType, code string) (*domain.NotificationChannelVerification, error) {
	verification := &domain.NotificationChannelVerification{}
	err := u.db.NewSelect().Model(verification).
		Where("ncv.channel = ? AND ncv.code = ? AND ncv.status = ?", channel, code, domain.ChannelVerificationPending).
		Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("channel verification not found: %v", err)
	}
	return verification, nil
}

func (u *DB) MarkChannelVerificationVerified(ctx context.Context, id, extendID string) error {
	now := time.Now()
	res, err := u.db.NewUpdate().Model((*domain.NotificationChannelVerification)(nil)).
		Set("status = ?", domain.ChannelVerificationVerified).
		Set("extend_id = ?", extendID).
		Set("verified_at = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("channel verification not updated: %v", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("channel verification not found")
	}
	return nil
}

func (u *DB) ExpirePendingChannelVerifications(ctx context.Context, userID string, channel domain.NotificationChannelType) error {
	_, err := u.db.NewUpdate().Model((*domain.NotificationChannelVerification)(nil)).
		Set("status = ?", domain.ChannelVerificationExpired).
		Set("updated_at = current_timestamp").
		Where("user_id = ? AND channel = ? AND status = ?", userID, channel, domain.ChannelVerificationPending).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("channel verifications not expired: %v", err)
	}
	return nil
}
