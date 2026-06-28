package services

import (
	"errors"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
)

type TaskService struct {
	repo ports.TaskRepository
}

func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(req domain.CreateTaskRequest) (*domain.Task, error) {
	countTask, err := s.repo.GetCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get task count: %v", err)
	}
	taskID := "TASK-" + fmt.Sprintf("%04d", countTask+1)

	status := domain.TaskStatusTodo
	if req.Status != "" {
		status = req.Status
	}

	priority := domain.TaskPriorityMedium
	if req.Priority != "" {
		priority = req.Priority
	}

	task := domain.Task{
		TaskID:      taskID,
		Title:       req.Title,
		Status:      status,
		Label:       req.Label,
		Priority:    priority,
		HTML:        req.HTML,
		Lexical:     req.Lexical,
		Description: req.Description,
	}

	created, err := s.repo.CreateTask(task)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *TaskService) GetTask(id string) (*domain.Task, error) {
	if id == "" {
		return nil, errors.New("task id is required")
	}

	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) UpdateTask(id string, req domain.UpdateTaskRequest) (*domain.Task, error) {
	if id == "" {
		return nil, errors.New("task id is required")
	}

	existing, err := s.repo.GetTaskByID(id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %v", err)
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Label != nil {
		existing.Label = *req.Label
	}
	if req.Priority != nil {
		existing.Priority = *req.Priority
	}
	if req.HTML != nil {
		existing.HTML = req.HTML
	}
	if req.Lexical != nil {
		existing.Lexical = req.Lexical
	}
	if req.Description != nil {
		existing.Description = req.Description
	}

	updated, err := s.repo.UpdateTask(id, *existing)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *TaskService) DeleteTask(id string) error {
	if id == "" {
		return errors.New("task id is required")
	}

	err := s.repo.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) ListTasks(filter domain.TaskFilter) ([]*domain.Task, int, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}

	tasks, total, err := s.repo.ListTasks(filter)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (s *TaskService) ListTasksCursor(filter domain.TaskFilter, cursor string, limit int) ([]*domain.Task, *string, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.repo.ListTasksCursor(filter, cursor, limit)
}

func (s *TaskService) ListAllTasks() ([]*domain.Task, error) {
	tasks, err := s.repo.ListAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetTaskStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	todoCount, _ := s.repo.CountTasksByStatus(domain.TaskStatusTodo)
	inProgressCount, _ := s.repo.CountTasksByStatus(domain.TaskStatusInProgress)
	doneCount, _ := s.repo.CountTasksByStatus(domain.TaskStatusDone)
	cancelledCount, _ := s.repo.CountTasksByStatus(domain.TaskStatusCancelled)

	stats["by_status"] = map[string]int{
		"todo":        todoCount,
		"in_progress": inProgressCount,
		"done":        doneCount,
		"cancelled":   cancelledCount,
	}

	highCount, _ := s.repo.CountTasksByPriority(domain.TaskPriorityHigh)
	criticalCount, _ := s.repo.CountTasksByPriority(domain.TaskPriorityCritical)

	stats["by_priority"] = map[string]int{
		"high":     highCount,
		"critical": criticalCount,
	}

	return stats, nil
}
