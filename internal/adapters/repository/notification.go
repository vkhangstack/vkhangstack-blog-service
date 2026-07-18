package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateNotification(ctx context.Context, notification domain.Notification) (*domain.Notification, error) {
	notification.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&notification).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("notification not saved: %v", err)
	}
	return &notification, nil
}

func (u *DB) GetNotificationByID(ctx context.Context, userID, id string) (*domain.Notification, error) {
	notification := &domain.Notification{}
	err := u.db.NewSelect().Model(notification).
		Where("n.id = ? AND n.user_id = ? AND n.is_deleted = FALSE", id, userID).
		Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("notification not found")
	}
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (u *DB) UpdateNotification(ctx context.Context, userID, id string, updates domain.Notification) error {
	updates.ID = id
	updates.UserID = userID
	res, err := u.db.NewUpdate().Model(&updates).OmitZero().
		Where("id = ? AND user_id = ?", id, userID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("notification not updated: %v", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("notification not found")
	}
	return nil
}

func (u *DB) DeleteNotification(ctx context.Context, userID, id string) error {
	res, err := u.db.NewUpdate().Model((*domain.Notification)(nil)).
		Set("is_deleted = TRUE").
		Where("id = ? AND user_id = ?", id, userID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("notification not deleted: %v", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("notification not found")
	}
	return nil
}

func (u *DB) ListNotifications(ctx context.Context, userID string, filter domain.NotificationFilter) ([]*domain.Notification, int, error) {
	notifications := make([]*domain.Notification, 0)
	query := u.db.NewSelect().Model(&notifications).
		Where("n.user_id = ? AND n.is_deleted = FALSE", userID)

	if filter.Type != "" {
		query = query.Where("n.type = ?", filter.Type)
	}
	if filter.Read != nil {
		query = query.Where("n.is_read = ?", *filter.Read)
	}

	total, err := query.Order("n.created_at DESC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("notifications not found: %v", err)
	}

	return notifications, total, nil
}

func (u *DB) ListNotificationsCursor(ctx context.Context, userID string, filter domain.NotificationFilter, cursor string, limit int) ([]*domain.Notification, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	var notifications []*domain.Notification
	query := u.db.NewSelect().Model(&notifications).
		Where("n.user_id = ? AND n.is_deleted = FALSE", userID)

	if filter.Type != "" {
		query = query.Where("n.type = ?", filter.Type)
	}
	if filter.Read != nil {
		query = query.Where("n.is_read = ?", *filter.Read)
	}

	queryCount := query.Clone()
	query = query.Order("n.created_at DESC", "n.id DESC")

	if cursorID != "" {
		cursorNotification := &domain.Notification{}
		err := u.db.NewSelect().Model(cursorNotification).
			Where("n.id = ? AND n.user_id = ?", cursorID, userID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor notification not found: %v", err)
		}
		query = query.Where("(n.created_at, n.id) < (?, ?)", cursorNotification.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("notifications not found: %v", err)
	}

	var nextCursor *string
	if len(notifications) > limit {
		notifications = notifications[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(notifications[len(notifications)-1].ID))
	}

	total, _ := queryCount.Count(ctx)
	return notifications, nextCursor, total, nil
}

func (u *DB) GetNotificationSettingByUserID(ctx context.Context, userID string) (*domain.NotificationSetting, error) {
	setting := &domain.NotificationSetting{}
	err := u.db.NewSelect().Model(setting).Where("ns.user_id = ?", userID).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("notification setting not found: %v", err)
	}
	return setting, nil
}

func (u *DB) UpsertNotificationSetting(ctx context.Context, setting domain.NotificationSetting) (*domain.NotificationSetting, error) {
	if setting.ID == "" {
		setting.ID = u.snowflakeNode.GenerateID()
	}
	_, err := u.db.NewInsert().Model(&setting).
		On("CONFLICT (user_id) DO UPDATE").
		Set("channels = EXCLUDED.channels").
		Set("updated_at = current_timestamp").
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("notification setting not saved: %v", err)
	}
	return u.GetNotificationSettingByUserID(ctx, setting.UserID)
}
