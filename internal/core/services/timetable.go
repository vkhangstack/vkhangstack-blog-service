package services

import (
	"context"
	"errors"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type TimetableService struct {
	repo ports.TimetableRepository
}

func NewTimetableService(repo ports.TimetableRepository) *TimetableService {
	return &TimetableService{repo: repo}
}

func (s *TimetableService) CreateTimetableEntry(ctx context.Context, authorID string, req domain.CreateTimetableEntryRequest) (*domain.TimetableEntry, error) {
	timetable := domain.TimetableEntry{
		AuthorID:  authorID,
		Subject:   req.Subject,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Color:     *req.Color,
		Note:      req.Note,
		CreatedBy: authorID,
	}
	created, err := s.repo.CreateTimetableEntry(ctx, timetable)
	if err != nil {
		return nil, err
	}

	return created, nil
}
func (s *TimetableService) GetTimetableEntry(ctx context.Context, authorID, timetableID string) (*domain.TimetableEntry, error) {
	timetable, err := s.repo.GetTimetableEntryByID(ctx, authorID, timetableID)
	if err != nil {
		return nil, err
	}

	return timetable, nil
}

func (s *TimetableService) ListTimetableEntries(ctx context.Context, authorID string, filter domain.TimetableEntryFilter) ([]*domain.TimetableEntryResponse, int, error) {
	timetables, total, err := s.repo.ListTimetableEntries(ctx, authorID, filter)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*domain.TimetableEntryResponse, 0, total)
	for _, timetable := range timetables {
		tm := &domain.TimetableEntryResponse{
			ID:        timetable.ID,
			Subject:   timetable.Subject,
			DayOfWeek: timetable.DayOfWeek,
			StartTime: timetable.StartTime,
			EndTime:   timetable.EndTime,
			Color:     timetable.Color,
			Note:      timetable.Note,
			CreatedAt: timetable.CreatedAt,
			UpdatedAt: timetable.UpdatedAt,
		}
		result = append(result, tm)
	}

	return result, total, nil
}

func (s *TimetableService) ListTimetableEntriesCursor(ctx context.Context, authorID string, filter domain.TimetableEntryFilter, cursor string, limit int) ([]*domain.TimetableEntryResponse, *string, int, error) {
	timetables, nextCursor, total, err := s.repo.ListTimetableEntriesCursor(ctx, authorID, filter, cursor, limit)
	if err != nil {
		return nil, nil, 0, err
	}

	result := make([]*domain.TimetableEntryResponse, 0, len(timetables))
	for _, timetable := range timetables {
		tm := &domain.TimetableEntryResponse{
			ID:        timetable.ID,
			Subject:   timetable.Subject,
			DayOfWeek: timetable.DayOfWeek,
			StartTime: timetable.StartTime,
			EndTime:   timetable.EndTime,
			Color:     timetable.Color,
			Note:      timetable.Note,
			CreatedAt: timetable.CreatedAt,
			UpdatedAt: timetable.UpdatedAt,
		}
		result = append(result, tm)
	}

	return result, nextCursor, total, nil
}

func (s *TimetableService) UpdateTimetableEntry(ctx context.Context, authorID, timetableID string, req domain.UpdateTimetableEntryRequest) (*domain.TimetableEntry, error) {
	timetable, err := s.repo.GetTimetableEntryByID(ctx, authorID, timetableID)
	if err != nil {
		return nil, err
	}
	utils.SetIfNotNil(&timetable.Subject, req.Subject)
	utils.SetIfNotNil(&timetable.DayOfWeek, req.DayOfWeek)
	utils.SetIfNotNil(&timetable.StartTime, req.StartTime)
	utils.SetIfNotNil(&timetable.EndTime, req.EndTime)
	utils.SetIfNotNil(&timetable.Color, req.Color)
	utils.SetIfNotNil(&timetable.Note, &req.Note)

	err = s.repo.UpdateTimetableEntry(ctx, timetableID, *timetable)
	if err != nil {
		return nil, err
	}

	return timetable, nil
}

func (s *TimetableService) DeleteTimetableEntry(ctx context.Context, authorID, timetableID string) error {
	timetable, err := s.repo.GetTimetableEntryByID(ctx, authorID, timetableID)
	if err != nil {
		return err
	}
	if timetable == nil {
		return errors.New("timetable entry not found")
	}

	err = s.repo.DeleteTimetableEntry(ctx, timetableID)
	if err != nil {
		return err
	}

	return nil
}
