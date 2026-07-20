package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateQuiz(ctx context.Context, quiz domain.Quiz, questions []domain.QuizQuestion) (*domain.Quiz, error) {
	quiz.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&quiz).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("quiz not saved: %v", err)
	}

	if err := u.insertQuizQuestions(ctx, quiz.ID, questions); err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (u *DB) GetQuizByID(ctx context.Context, id string) (*domain.QuizWithQuestions, error) {
	quiz := &domain.Quiz{}
	err := u.db.NewSelect().Model(quiz).Where("q.id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("quiz not found: %v", err)
	}

	var questions []*domain.QuizQuestion
	err = u.db.NewSelect().
		Model(&questions).
		Where("qq.quiz_id = ?", id).
		Order("qq.position ASC", "qq.id ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("quiz questions not found: %v", err)
	}

	return &domain.QuizWithQuestions{Quiz: quiz, Questions: questions}, nil
}

func (u *DB) UpdateQuiz(ctx context.Context, id string, updates domain.Quiz) error {
	updates.ID = id
	_, err := u.db.NewUpdate().Model(&updates).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("quiz not updated: %v", err)
	}
	return nil
}

func (u *DB) ReplaceQuizQuestions(ctx context.Context, quizID string, questions []domain.QuizQuestion) error {
	_, err := u.db.NewDelete().Model((*domain.QuizQuestion)(nil)).Where("quiz_id = ?", quizID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("quiz questions not replaced: %v", err)
	}
	return u.insertQuizQuestions(ctx, quizID, questions)
}

func (u *DB) DeleteQuiz(ctx context.Context, id string) error {
	quiz := &domain.Quiz{ID: id}
	_, err := u.db.NewDelete().Model(quiz).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("quiz not deleted: %v", err)
	}
	_, err = u.db.NewDelete().Model((*domain.QuizQuestion)(nil)).Where("quiz_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("quiz questions not deleted: %v", err)
	}
	return nil
}

func (u *DB) ListQuizzesCursor(ctx context.Context, filter domain.QuizFilter, cursor string, limit int) ([]*domain.QuizWithCount, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	var quizzes []*domain.Quiz
	query := u.db.NewSelect().Model(&quizzes)

	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("q.status = ?", *filter.Status)
	}
	if filter.ExamType != nil && *filter.ExamType != "" {
		query = query.Where("q.exam_type = ?", *filter.ExamType)
	}
	if filter.Skill != nil && *filter.Skill != "" {
		query = query.Where("q.skill = ?", *filter.Skill)
	}
	if filter.CreatedBy != nil && *filter.CreatedBy != "" {
		query = query.Where("q.created_by = ?", *filter.CreatedBy)
	}
	if filter.Title != nil && *filter.Title != "" {
		query = query.Where("q.title ILIKE ?", "%"+*filter.Title+"%")
	}

	queryCount := query.Clone()
	query = query.Order("q.created_at DESC", "q.id DESC")

	if cursorID != "" {
		cursorQuiz := &domain.Quiz{}
		err := u.db.NewSelect().Model(cursorQuiz).Where("q.id = ?", cursorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor quiz not found: %v", err)
		}
		query = query.Where("(q.created_at, q.id) < (?, ?)", cursorQuiz.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("quizzes not found: %v", err)
	}

	var nextCursor *string
	if len(quizzes) > limit {
		quizzes = quizzes[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(quizzes[len(quizzes)-1].ID))
	}
	total, _ := queryCount.Count(ctx)

	var quizIDs []string
	for _, quiz := range quizzes {
		quizIDs = append(quizIDs, quiz.ID)
	}

	counts, err := u.getQuestionCounts(ctx, quizIDs)
	if err != nil {
		return nil, nil, 0, err
	}

	result := []*domain.QuizWithCount{}
	for _, quiz := range quizzes {
		result = append(result, &domain.QuizWithCount{
			Quiz:          quiz,
			QuestionCount: counts[quiz.ID],
		})
	}
	return result, nextCursor, total, nil
}

func (u *DB) CreateQuizAttempt(ctx context.Context, attempt domain.QuizAttempt) (*domain.QuizAttempt, error) {
	attempt.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&attempt).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("quiz attempt not saved: %v", err)
	}
	return &attempt, nil
}

func (u *DB) ListQuizAttemptsCursor(ctx context.Context, quizID string, cursor string, limit int) ([]*domain.QuizAttempt, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	var attempts []*domain.QuizAttempt
	query := u.db.NewSelect().Model(&attempts).Where("qa.quiz_id = ?", quizID)

	queryCount := query.Clone()
	query = query.Order("qa.created_at DESC", "qa.id DESC")

	if cursorID != "" {
		cursorAttempt := &domain.QuizAttempt{}
		err := u.db.NewSelect().Model(cursorAttempt).Where("qa.id = ?", cursorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor attempt not found: %v", err)
		}
		query = query.Where("(qa.created_at, qa.id) < (?, ?)", cursorAttempt.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("quiz attempts not found: %v", err)
	}

	var nextCursor *string
	if len(attempts) > limit {
		attempts = attempts[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(attempts[len(attempts)-1].ID))
	}
	total, _ := queryCount.Count(ctx)

	return attempts, nextCursor, total, nil
}

func (u *DB) insertQuizQuestions(ctx context.Context, quizID string, questions []domain.QuizQuestion) error {
	if len(questions) == 0 {
		return nil
	}
	for i := range questions {
		questions[i].ID = u.snowflakeNode.GenerateID()
		questions[i].QuizID = quizID
		questions[i].Position = i
	}
	_, err := u.db.NewInsert().Model(&questions).Exec(ctx)
	if err != nil {
		return fmt.Errorf("quiz questions not saved: %v", err)
	}
	return nil
}

func (u *DB) getQuestionCounts(ctx context.Context, quizIDs []string) (map[string]int, error) {
	if len(quizIDs) == 0 {
		return map[string]int{}, nil
	}

	var rows []struct {
		QuizID string `bun:"quiz_id"`
		Count  int    `bun:"count"`
	}

	err := u.db.NewSelect().
		Model((*domain.QuizQuestion)(nil)).
		ColumnExpr("quiz_id, COUNT(*) AS count").
		Where("quiz_id IN (?)", bun.In(quizIDs)).
		Group("quiz_id").
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get question counts: %w", err)
	}

	counts := make(map[string]int, len(quizIDs))
	for _, r := range rows {
		counts[r.QuizID] = r.Count
	}
	return counts, nil
}
