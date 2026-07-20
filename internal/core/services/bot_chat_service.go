package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
)

// BotChatService handles plain-text bot commands sent via the Zalo bot channel.
// Supported commands: /tasks, /notes, /timetable
type BotChatService struct {
	settingRepo  ports.NotificationSettingRepository
	taskSvc      *TaskService
	noteSvc      *NoteService
	timetableSvc *TimetableService
	bot          ports.ZaloBotClient
	chatSvc      *ChatService
}

func NewBotChatService(
	settingRepo ports.NotificationSettingRepository,
	taskSvc *TaskService,
	noteSvc *NoteService,
	timetableSvc *TimetableService,
	bot ports.ZaloBotClient,
	chatSvc *ChatService,
) *BotChatService {
	return &BotChatService{
		settingRepo:  settingRepo,
		taskSvc:      taskSvc,
		noteSvc:      noteSvc,
		timetableSvc: timetableSvc,
		bot:          bot,
		chatSvc:      chatSvc,
	}
}

// HandleMessage resolves the sender by their Zalo extend_id, then dispatches
// the command and replies with a formatted list.
//
// Supported syntax:
//
//	/tasks [status]   — list tasks, optionally filtered by status (todo|in_progress|done|cancelled)
//	/notes [status]   — list notes, optionally filtered by status
func (s *BotChatService) HandleMessage(ctx context.Context, senderID, text string) {
	parts := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	if len(parts) == 0 {
		s.reply(senderID, "Commands:\n/tasks [status]\n/notes [status]\n/timetable [day 1-7]")
		return
	}

	cmd := parts[0]
	var filter string
	if len(parts) > 1 {
		filter = parts[1]
	}

	if cmd != "/tasks" && cmd != "/notes" && cmd != "/timetable" {
		s.reply(senderID, "Commands:\n/tasks [status]\n/notes [status]\n/timetable [day 1-7], /timetable [day 1-7]")
		return
	}

	// Look up the user who owns this Zalo extend_id.
	setting, err := s.settingRepo.GetNotificationSettingByChannelToken(ctx, domain.NotificationChannelZaloBot, senderID)
	if err != nil {
		logger.Log.WithError(err).Error("BotChatService: failed to find user by channel token")
		s.reply(senderID, "Could not find your account. Please link your Zalo account first.")
		return
	}
	if setting == nil {
		s.reply(senderID, "Your Zalo account is not linked. Please link it first via the app.")
		return
	}

	switch cmd {
	case "/tasks":
		s.handleTasks(ctx, senderID, setting.UserID, filter)
	case "/notes":
		s.handleNotes(ctx, senderID, setting.UserID, filter)
	case "/timetable":
		s.handleTimetable(ctx, senderID, setting.UserID, filter)
	}
}

func (s *BotChatService) handleTasks(ctx context.Context, senderID, userID, statusFilter string) {
	taskFilter := domain.TaskFilter{
		CreatedBy: userID,
		Page:      1,
		Limit:     10,
	}
	if statusFilter != "" {
		taskFilter.Status = statusFilter
	}

	tasks, _, err := s.taskSvc.ListTasks(ctx, taskFilter)
	if err != nil {
		logger.Log.WithError(err).Error("BotChatService: failed to list tasks")
		s.reply(senderID, "Failed to fetch your tasks.")
		return
	}
	if len(tasks) == 0 {
		label := "tasks"
		if statusFilter != "" {
			label = fmt.Sprintf("%s tasks", statusFilter)
		}
		s.reply(senderID, fmt.Sprintf("You have no %s.", label))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Your tasks (%d):\n", len(tasks)))
	for i, t := range tasks {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s (%s)\n", i+1, strings.ToUpper(string(t.Status)), t.Title, t.Priority))
	}
	s.reply(senderID, sb.String())
}

func (s *BotChatService) handleNotes(ctx context.Context, senderID, userID, statusFilter string) {
	noteFilter := domain.NoteFilter{
		CreatedBy: &userID,
		Page:      1,
		Limit:     10,
	}
	if statusFilter != "" {
		noteFilter.Status = &statusFilter
	}

	notes, _, err := s.noteSvc.ListNotes(ctx, noteFilter)
	if err != nil {
		logger.Log.WithError(err).Error("BotChatService: failed to list notes")
		s.reply(senderID, "Failed to fetch your notes.")
		return
	}
	if len(notes) == 0 {
		label := "notes"
		if statusFilter != "" {
			label = fmt.Sprintf("%s notes", statusFilter)
		}
		s.reply(senderID, fmt.Sprintf("You have no %s.", label))
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Your notes (%d):\n", len(notes)))
	for i, n := range notes {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, n.Status, n.Title))
	}
	s.reply(senderID, sb.String())
}

// handleTimetable lists the user's timetable entries, optionally filtered by day of week (1–7).
func (s *BotChatService) handleTimetable(ctx context.Context, senderID, userID, dayFilter string) {
	entryFilter := domain.TimetableEntryFilter{Page: 1, Limit: 20}
	if dayFilter != "" {
		day, err := strconv.Atoi(dayFilter)
		if err != nil || day < 1 || day > 7 {
			s.reply(senderID, "Invalid day. Use a number 1 (Mon) – 7 (Sun).")
			return
		}
		entryFilter.DayOfWeek = &day
	}

	entries, _, err := s.timetableSvc.ListTimetableEntries(ctx, userID, entryFilter)
	if err != nil {
		logger.Log.WithError(err).Error("BotChatService: failed to list timetable entries")
		s.reply(senderID, "Failed to fetch your timetable.")
		return
	}
	if len(entries) == 0 {
		msg := "You have no timetable entries."
		if dayFilter != "" {
			msg = fmt.Sprintf("No timetable entries for day %s.", dayFilter)
		}
		s.reply(senderID, msg)
		return
	}

	dayNames := [...]string{"", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Your timetable (%d):\n", len(entries)))
	for i, e := range entries {
		day := ""
		if e.DayOfWeek >= 1 && e.DayOfWeek <= 7 {
			day = dayNames[e.DayOfWeek]
		}
		sb.WriteString(fmt.Sprintf("%d. [%s] %s %s–%s\n", i+1, day, e.Subject, e.StartTime, e.EndTime))
	}
	s.reply(senderID, sb.String())
}

func (s *BotChatService) reply(senderID, text string) {
	if err := s.bot.SendMessage(senderID, text); err != nil {
		logger.Log.WithError(err).WithField("sender_id", senderID).Error("BotChatService: failed to send reply")
		return
	}
	if s.chatSvc != nil {
		s.chatSvc.RecordOutbound(context.Background(), senderID, domain.BotChatSenderBot, text)
	}
}
