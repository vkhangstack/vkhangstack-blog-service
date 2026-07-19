package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type NotificationService struct {
	repo ports.NotificationRepository
}

func NewNotificationService(repo ports.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func toNotificationResponse(n *domain.Notification) *domain.NotificationResponse {
	return &domain.NotificationResponse{
		ID:        n.ID,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type,
		Read:      n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID string, req domain.CreateNotificationRequest) (*domain.NotificationResponse, error) {
	notification := domain.Notification{
		UserID:  userID,
		Type:    req.Type,
		Title:   req.Title,
		Message: req.Message,
	}

	created, err := s.repo.CreateNotification(ctx, notification)
	if err != nil {
		return nil, err
	}
	return toNotificationResponse(created), nil
}

func (s *NotificationService) GetNotification(ctx context.Context, userID, id string) (*domain.NotificationResponse, error) {
	if id == "" {
		return nil, errors.New("notification id is required")
	}
	notification, err := s.repo.GetNotificationByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return toNotificationResponse(notification), nil
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID string, filter domain.NotificationFilter) ([]*domain.NotificationResponse, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}

	notifications, total, err := s.repo.ListNotifications(ctx, userID, filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*domain.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toNotificationResponse(n))
	}
	return result, total, nil
}

func (s *NotificationService) ListNotificationsCursor(ctx context.Context, userID string, filter domain.NotificationFilter, cursor string, limit int) ([]*domain.NotificationResponse, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	notifications, nextCursor, total, err := s.repo.ListNotificationsCursor(ctx, userID, filter, cursor, limit)
	if err != nil {
		return nil, nil, 0, err
	}

	result := make([]*domain.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, toNotificationResponse(n))
	}
	return result, nextCursor, total, nil
}

func (s *NotificationService) UpdateNotification(ctx context.Context, userID, id string, req domain.UpdateNotificationRequest) error {
	if id == "" {
		return errors.New("notification id is required")
	}

	updates := domain.Notification{}
	if req.Read != nil {
		updates.IsRead = *req.Read
	}

	return s.repo.UpdateNotification(ctx, userID, id, updates)
}

func (s *NotificationService) DeleteNotification(ctx context.Context, userID, id string) error {
	if id == "" {
		return errors.New("notification id is required")
	}
	return s.repo.DeleteNotification(ctx, userID, id)
}

// NotificationSettingService manages a user's per-channel notification preferences
// (which channels are enabled and the token/address used to deliver to each one), including
// channels that must be linked by having the user copy a code into the channel's own app.
type NotificationSettingService struct {
	repo             ports.NotificationSettingRepository
	verificationRepo ports.NotificationChannelVerificationRepository
	zaloBot          ports.ZaloBotClient
}

func NewNotificationSettingService(repo ports.NotificationSettingRepository, verificationRepo ports.NotificationChannelVerificationRepository, zaloBot ports.ZaloBotClient) *NotificationSettingService {
	return &NotificationSettingService{repo: repo, verificationRepo: verificationRepo, zaloBot: zaloBot}
}

// VerificationCodePattern extracts a linking code (e.g. "NEXION-HUB-482913") from arbitrary
// message text a user sends to a bot channel.
var VerificationCodePattern = regexp.MustCompile(`NEXION-HUB-\d{6}`)

const verificationCodeTTL = 15 * time.Minute

func generateVerificationCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return "NEXION-HUB-" + string(code), nil
}

func toNotificationSettingResponse(s *domain.NotificationSetting) *domain.NotificationSettingResponse {
	return &domain.NotificationSettingResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		Channels:  s.Channels,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// GetNotificationSetting returns the user's settings, or an empty (no channels enabled)
// response if they haven't configured any yet — no row is created until they save one.
func (s *NotificationSettingService) GetNotificationSetting(ctx context.Context, userID string) (*domain.NotificationSettingResponse, error) {
	setting, err := s.repo.GetNotificationSettingByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return &domain.NotificationSettingResponse{
			UserID:   userID,
			Channels: []domain.NotificationChannel{},
		}, nil
	}
	return toNotificationSettingResponse(setting), nil
}

func (s *NotificationSettingService) UpdateNotificationSetting(ctx context.Context, userID string, req domain.UpdateNotificationSettingRequest) (*domain.NotificationSettingResponse, error) {
	channels := make([]domain.NotificationChannel, 0, len(req.Channels))
	seen := make(map[domain.NotificationChannelType]bool, len(req.Channels))
	for _, c := range req.Channels {
		if !domain.IsKnownNotificationChannel(c.Channel) {
			return nil, errors.New("unknown notification channel: " + string(c.Channel))
		}
		if seen[c.Channel] {
			return nil, errors.New("duplicate notification channel: " + string(c.Channel))
		}
		seen[c.Channel] = true
		channels = append(channels, domain.NotificationChannel{
			Channel: c.Channel,
			Enabled: c.Enabled,
			Token:   c.Token,
		})
	}

	setting := domain.NotificationSetting{
		UserID:   userID,
		Channels: channels,
	}

	saved, err := s.repo.UpsertNotificationSetting(ctx, setting)
	if err != nil {
		return nil, err
	}
	return toNotificationSettingResponse(saved), nil
}

// RequestChannelVerification generates a code the user copies into the channel's own app
// (e.g. sends it as a Zalo message to our bot) to link that channel without a client-supplied token.
func (s *NotificationSettingService) RequestChannelVerification(ctx context.Context, userID string, req domain.CreateChannelVerificationRequest) (*domain.ChannelVerificationResponse, error) {
	if !domain.IsVerifiableNotificationChannel(req.Channel) {
		return nil, errors.New("channel does not support verification linking: " + string(req.Channel))
	}

	if err := s.verificationRepo.ExpirePendingChannelVerifications(ctx, userID, req.Channel); err != nil {
		return nil, err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(verificationCodeTTL)
	verification := domain.NotificationChannelVerification{
		UserID:    userID,
		Channel:   req.Channel,
		Code:      code,
		Status:    domain.ChannelVerificationPending,
		ExpiresAt: expiresAt,
	}
	if _, err := s.verificationRepo.CreateChannelVerification(ctx, verification); err != nil {
		return nil, err
	}

	return &domain.ChannelVerificationResponse{
		Channel:   req.Channel,
		Code:      code,
		Message:   fmt.Sprintf("Send this code to our Zalo bot to link your account: %s", code),
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyChannelNotificationLinked is called by a channel's webhook handler once it receives a
// message carrying a verification code; on success it enables that channel with the
// sender's platform-specific ID (extendID) as the channel's token.
func (s *NotificationSettingService) VerifyChannelNotificationLinked(ctx context.Context, channel domain.NotificationChannelType, code, extendID string) error {
	matchedCode := VerificationCodePattern.FindString(code)
	if matchedCode == "" {
		return errors.New("no verification code found in message")
	}

	verification, err := s.verificationRepo.GetPendingChannelVerification(ctx, channel, matchedCode)
	if err != nil {
		return err
	}
	if verification == nil || time.Now().After(verification.ExpiresAt) {
		return errors.New("verification code invalid or expired")
	}

	if err := s.verificationRepo.MarkChannelVerificationVerified(ctx, verification.ID, extendID); err != nil {
		return err
	}

	setting, err := s.repo.GetNotificationSettingByUserID(ctx, verification.UserID)
	if err != nil {
		return err
	}

	channels := []domain.NotificationChannel{}
	if setting != nil {
		channels = setting.Channels
	}
	found := false
	for i, c := range channels {
		if c.Channel == channel {
			channels[i].Enabled = true
			channels[i].Token = extendID
			found = true
			break
		}
	}
	if !found {
		channels = append(channels, domain.NotificationChannel{
			Channel: channel,
			Enabled: true,
			Token:   extendID,
		})
	}

	if _, err := s.repo.UpsertNotificationSetting(ctx, domain.NotificationSetting{
		UserID:   verification.UserID,
		Channels: channels,
	}); err != nil {
		return err
	}

	if channel == domain.NotificationChannelZaloBot && s.zaloBot != nil {
		_ = s.zaloBot.SendMessage(extendID, "Your Zalo channel has been linked successfully. You'll now receive notifications here.")
	}

	return nil
}

// UnlinkChannel disables a linked channel and clears its token, so the user can restart the
// linking flow (e.g. RequestChannelVerification) from a clean state.
func (s *NotificationSettingService) UnlinkChannel(ctx context.Context, userID string, channel domain.NotificationChannelType) (*domain.NotificationSettingResponse, error) {
	if !domain.IsKnownNotificationChannel(channel) {
		return nil, errors.New("unknown notification channel: " + string(channel))
	}

	setting, err := s.repo.GetNotificationSettingByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return &domain.NotificationSettingResponse{
			UserID:   userID,
			Channels: []domain.NotificationChannel{},
		}, nil
	}

	var unlinkedToken string
	channels := make([]domain.NotificationChannel, len(setting.Channels))
	copy(channels, setting.Channels)
	for i, c := range channels {
		if c.Channel == channel {
			unlinkedToken = c.Token
			channels[i].Enabled = false
			channels[i].Token = ""
		}
	}

	if err := s.verificationRepo.ExpirePendingChannelVerifications(ctx, userID, channel); err != nil {
		return nil, err
	}

	saved, err := s.repo.UpsertNotificationSetting(ctx, domain.NotificationSetting{
		UserID:   userID,
		Channels: channels,
	})
	if err != nil {
		return nil, err
	}

	if channel == domain.NotificationChannelZaloBot && unlinkedToken != "" && s.zaloBot != nil {
		_ = s.zaloBot.SendMessage(unlinkedToken, "Your Zalo channel has been unlinked. You'll no longer receive notifications here.")
	}

	return toNotificationSettingResponse(saved), nil
}
