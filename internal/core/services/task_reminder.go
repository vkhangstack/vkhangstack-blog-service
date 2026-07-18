package services

import (
	"context"
	"fmt"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// TaskReminderService checks tasks whose reminder_at has passed and, for those with notice
// enabled, delivers a notification to the assignee (in-app, plus their enabled channels).
type TaskReminderService struct {
	taskRepo     ports.TaskRepository
	notification ports.NotificationRepository
	setting      ports.NotificationSettingRepository
	zaloBot      ports.ZaloBotClient
}

func NewTaskReminderService(taskRepo ports.TaskRepository, notification ports.NotificationRepository, setting ports.NotificationSettingRepository, zaloBot ports.ZaloBotClient) *TaskReminderService {
	return &TaskReminderService{taskRepo: taskRepo, notification: notification, setting: setting, zaloBot: zaloBot}
}

// SendDueReminders finds tasks with enable_notice=true whose reminder_at has passed, delivers
// a notification to each task's assignee, then clears reminder_at so it isn't resent.
func (s *TaskReminderService) SendDueReminders(ctx context.Context) error {
	tasks, err := s.taskRepo.ListDueReminderTasks(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to list due reminder tasks: %w", err)
	}

	for _, task := range tasks {
		fmt.Printf("Sending reminder for task %s to assignee %v\n", task.ID, task.AssigneeID)
		if task.AssigneeID == nil || *task.AssigneeID == "" {
			s.clearReminder(ctx, task)
			continue
		}

		if err := s.notify(ctx, *task.AssigneeID, task); err != nil {
			logger.Log.WithError(err).WithField("task_id", task.ID).Error("TaskReminderService: failed to notify assignee")
			continue
		}

		s.clearReminder(ctx, task)
	}

	return nil
}

func (s *TaskReminderService) notify(ctx context.Context, userID string, task *domain.Task) error {
	message := fmt.Sprintf("Task \"%s\" is due for a reminder.", task.Title)

	if _, err := s.notification.CreateNotification(ctx, domain.Notification{
		UserID:  userID,
		Type:    "task_reminder",
		Title:   "Task reminder",
		Message: message,
	}); err != nil {
		return fmt.Errorf("failed to create in-app notification: %w", err)
	}

	setting, err := s.setting.GetNotificationSettingByUserID(ctx, userID)
	if err != nil || setting == nil {
		return nil
	}

	for _, channel := range setting.Channels {
		if !channel.Enabled || channel.Token == "" {
			continue
		}
		if channel.Channel == domain.NotificationChannelZaloBot && s.zaloBot != nil {
			if err := s.zaloBot.SendMessage(channel.Token, message); err != nil {
				fmt.Printf("err %s", err.Error())
				// logger.Log.WithError(err).WithField("task_id", task.ID).Error("TaskReminderService: failed to send zalo bot message")
			}
		}
	}

	return nil
}

func (s *TaskReminderService) clearReminder(ctx context.Context, task *domain.Task) {
	task.ReminderAt = nil
	if _, err := s.taskRepo.UpdateTask(ctx, task.ID, *task); err != nil {
		logger.Log.WithError(err).WithField("task_id", task.ID).Error("TaskReminderService: failed to clear reminder_at")
	}
}

// StartReminderPoller runs SendDueReminders on a fixed interval until ctx is cancelled.
func (s *TaskReminderService) StartReminderPoller(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.SendDueReminders(ctx); err != nil {
					logger.Log.WithError(err).Error("TaskReminderService: reminder poll failed")
				}
			}
		}
	}()
}
