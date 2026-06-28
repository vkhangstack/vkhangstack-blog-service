package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateTask(task domain.Task) (*domain.Task, error) {
	ctx := context.Background()
	task.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&task).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("task not saved: %v", err)
	}
	return &task, nil
}

func (u *DB) GetTaskByID(id string) (*domain.Task, error) {
	ctx := context.Background()
	task := &domain.Task{}
	err := u.db.NewSelect().Model(task).Where("t.id = ?", id).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	return task, err
}

func (u *DB) GetTaskByTaskID(taskID string) (*domain.Task, error) {
	ctx := context.Background()
	task := &domain.Task{}
	err := u.db.NewSelect().Model(task).Where("t.task_id = ?", taskID).Limit(1).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}
	return task, err
}

func (u *DB) UpdateTask(id string, updates domain.Task) (*domain.Task, error) {
	ctx := context.Background()
	updates.ID = id
	_, err := u.db.NewUpdate().Model(&updates).WherePK().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("task not updated: %v", err)
	}
	return u.GetTaskByID(id)
}

func (u *DB) DeleteTask(id string) error {
	ctx := context.Background()
	task := &domain.Task{ID: id}
	_, err := u.db.NewDelete().Model(task).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("task not deleted: %v", err)
	}
	return nil
}

func (u *DB) ListTasks(filter domain.TaskFilter) ([]*domain.Task, int, error) {
	ctx := context.Background()
	tasks := make([]*domain.Task, 0)
	query := u.db.NewSelect().Model(&tasks)

	if filter.Status != "" {
		query = query.Where("t.status = ?", filter.Status)
	}
	if filter.Label != "" {
		query = query.Where("t.label = ?", filter.Label)
	}
	if filter.Priority != "" {
		query = query.Where("t.priority = ?", filter.Priority)
	}

	total, err := query.Order("t.created_at DESC").Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("tasks not found: %v", err)
	}

	return tasks, total, nil
}

// ListTasksCursor returns tasks using cursor-based pagination
func (u *DB) ListTasksCursor(filter domain.TaskFilter, cursor string, limit int) ([]*domain.Task, *string, error) {
	ctx := context.Background()

	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, err
		}
		cursorID = id
	}

	var tasks []*domain.Task
	query := u.db.NewSelect().Model(&tasks)

	if filter.Status != "" {
		query = query.Where("t.status = ?", filter.Status)
	}
	if filter.Label != "" {
		query = query.Where("t.label = ?", filter.Label)
	}
	if filter.Priority != "" {
		query = query.Where("t.priority = ?", filter.Priority)
	}

	query = query.Order("t.created_at DESC", "t.id DESC")

	if cursorID != "" {
		cursorTask := &domain.Task{}
		err := u.db.NewSelect().Model(cursorTask).Where("t.id = ?", cursorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("cursor task not found: %v", err)
		}
		query = query.Where("(t.created_at, t.id) < (?, ?)", cursorTask.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("tasks not found: %v", err)
	}

	var nextCursor *string
	if len(tasks) > limit {
		tasks = tasks[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(tasks[len(tasks)-1].ID))
	}

	return tasks, nextCursor, nil
}

func (u *DB) ListAllTasks() ([]*domain.Task, error) {
	ctx := context.Background()
	tasks := make([]*domain.Task, 0)
	err := u.db.NewSelect().Model(&tasks).Order("t.created_at DESC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("tasks not found: %v", err)
	}
	return tasks, nil
}

func (u *DB) CountTasksByStatus(status domain.TaskStatus) (int, error) {
	ctx := context.Background()
	count, err := u.db.NewSelect().Model((*domain.Task)(nil)).Where("status = ?", status).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %v", err)
	}
	return count, nil
}

func (u *DB) CountTasksByPriority(priority domain.TaskPriority) (int, error) {
	ctx := context.Background()
	count, err := u.db.NewSelect().Model((*domain.Task)(nil)).Where("priority = ?", priority).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %v", err)
	}
	return count, nil
}

func (u *DB) GetCount() (int, error) {
	ctx := context.Background()
	count, err := u.db.NewSelect().Model((*domain.Task)(nil)).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %v", err)
	}
	return count, nil
}
