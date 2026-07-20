package services

import (
	"context"
	"fmt"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type QuizService struct {
	repo      ports.QuizRepository
	generator ports.QuizGenerator
}

func NewQuizService(repo ports.QuizRepository, generator ports.QuizGenerator) *QuizService {
	return &QuizService{
		repo:      repo,
		generator: generator,
	}
}

// GenerateQuiz drafts a quiz via the configured AI provider without persisting it.
func (s *QuizService) GenerateQuiz(ctx context.Context, req domain.GenerateQuizRequest) (*domain.CreateQuizRequest, error) {
	if s.generator == nil {
		return nil, fmt.Errorf("AI quiz generation is not configured (set AI_PROVIDER and AI_API_KEY)")
	}
	draft, err := s.generator.GenerateQuiz(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := validateQuizQuestions(draft.Questions); err != nil {
		return nil, err
	}
	return draft, nil
}

func validateQuizQuestions(questions []domain.QuizQuestionInput) error {
	for i, q := range questions {
		if q.CorrectIndex < 0 || q.CorrectIndex >= len(q.Options) {
			return fmt.Errorf("question %d: correct_index %d out of range for %d options", i+1, q.CorrectIndex, len(q.Options))
		}
	}
	return nil
}

func mapQuizQuestionInputs(inputs []domain.QuizQuestionInput) []domain.QuizQuestion {
	questions := make([]domain.QuizQuestion, 0, len(inputs))
	for _, in := range inputs {
		questions = append(questions, domain.QuizQuestion{
			Prompt:       in.Prompt,
			Options:      in.Options,
			CorrectIndex: in.CorrectIndex,
			Explanation:  in.Explanation,
		})
	}
	return questions
}

func (s *QuizService) CreateQuiz(ctx context.Context, authorID string, req domain.CreateQuizRequest) (*domain.Quiz, error) {
	if err := validateQuizQuestions(req.Questions); err != nil {
		return nil, err
	}

	examType := req.ExamType
	if examType == "" {
		examType = domain.QuizExamTypeGeneral
	}
	status := req.Status
	if status == "" {
		status = domain.QuizStatusDraft
	}

	quiz := domain.Quiz{
		Title:           req.Title,
		Description:     req.Description,
		ExamType:        examType,
		Skill:           req.Skill,
		Status:          status,
		DurationMinutes: req.DurationMinutes,
		CreatedBy:       utils.StringPtr(authorID),
	}

	return s.repo.CreateQuiz(ctx, quiz, mapQuizQuestionInputs(req.Questions))
}

func (s *QuizService) GetQuiz(ctx context.Context, id string) (*domain.QuizWithQuestions, error) {
	return s.repo.GetQuizByID(ctx, id)
}

func (s *QuizService) ListQuizzesCursor(ctx context.Context, filter domain.QuizFilter, cursor string, limit int) ([]*domain.QuizWithCount, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return s.repo.ListQuizzesCursor(ctx, filter, cursor, limit)
}

func (s *QuizService) UpdateQuiz(ctx context.Context, id string, req domain.UpdateQuizRequest) error {
	quizInfo, err := s.repo.GetQuizByID(ctx, id)
	if err != nil {
		return err
	}

	quiz := *quizInfo.Quiz
	utils.SetIfNotNil(&quiz.Title, req.Title)
	utils.SetIfNotNil(&quiz.ExamType, req.ExamType)
	utils.SetIfNotNil(&quiz.Status, req.Status)
	if req.Description != nil {
		quiz.Description = req.Description
	}
	if req.Skill != nil {
		quiz.Skill = req.Skill
	}
	if req.DurationMinutes != nil {
		quiz.DurationMinutes = req.DurationMinutes
	}
	quiz.UpdatedBy = utils.StringPtr(ctx.Value("user_id").(string))

	if err := s.repo.UpdateQuiz(ctx, id, quiz); err != nil {
		return err
	}

	if req.Questions != nil {
		if err := validateQuizQuestions(*req.Questions); err != nil {
			return err
		}
		return s.repo.ReplaceQuizQuestions(ctx, id, mapQuizQuestionInputs(*req.Questions))
	}
	return nil
}

func (s *QuizService) DeleteQuiz(ctx context.Context, id string) error {
	return s.repo.DeleteQuiz(ctx, id)
}

// SubmitAttempt grades the given answers server-side against the quiz's questions,
// persists the attempt, and returns per-question results with explanations.
func (s *QuizService) SubmitAttempt(ctx context.Context, quizID string, authorID string, req domain.SubmitQuizAttemptRequest) (*domain.QuizAttemptResult, error) {
	quiz, err := s.repo.GetQuizByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if len(quiz.Questions) == 0 {
		return nil, fmt.Errorf("quiz has no questions")
	}

	score := 0
	results := make([]domain.QuizQuestionResult, 0, len(quiz.Questions))
	for _, question := range quiz.Questions {
		selected, answered := req.Answers[question.ID]
		if !answered {
			selected = -1
		}
		correct := answered && selected == question.CorrectIndex
		if correct {
			score++
		}
		results = append(results, domain.QuizQuestionResult{
			QuestionID:    question.ID,
			SelectedIndex: selected,
			CorrectIndex:  question.CorrectIndex,
			Correct:       correct,
			Explanation:   question.Explanation,
		})
	}

	attempt := domain.QuizAttempt{
		QuizID:          quizID,
		Answers:         req.Answers,
		Score:           score,
		Total:           len(quiz.Questions),
		DurationSeconds: req.DurationSeconds,
		CreatedBy:       utils.StringPtr(authorID),
	}
	saved, err := s.repo.CreateQuizAttempt(ctx, attempt)
	if err != nil {
		return nil, err
	}

	return &domain.QuizAttemptResult{
		AttemptID: saved.ID,
		Score:     score,
		Total:     len(quiz.Questions),
		Results:   results,
	}, nil
}

func (s *QuizService) ListQuizAttemptsCursor(ctx context.Context, quizID string, cursor string, limit int) ([]*domain.QuizAttempt, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return s.repo.ListQuizAttemptsCursor(ctx, quizID, cursor, limit)
}
