package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/ports"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

const (
	defaultEaseFactor = 2.5
	minEaseFactor     = 1.3
)

type FlashcardService struct {
	repo      ports.FlashcardRepository
	generator ports.FlashcardGenerator
}

func NewFlashcardService(repo ports.FlashcardRepository, generator ports.FlashcardGenerator) *FlashcardService {
	return &FlashcardService{
		repo:      repo,
		generator: generator,
	}
}

func mapFlashcardCardInputs(inputs []domain.FlashcardCardInput) []domain.FlashcardCard {
	cards := make([]domain.FlashcardCard, 0, len(inputs))
	for _, in := range inputs {
		cards = append(cards, domain.FlashcardCard{
			Front:    in.Front,
			Back:     in.Back,
			Example:  in.Example,
			Phonetic: in.Phonetic,
		})
	}
	return cards
}

func (s *FlashcardService) CreateDeck(ctx context.Context, authorID string, req domain.CreateFlashcardDeckRequest) (*domain.FlashcardDeck, error) {
	language := req.Language
	if language == "" {
		language = "en"
	}
	status := req.Status
	if status == "" {
		status = domain.FlashcardDeckStatusDraft
	}

	deck := domain.FlashcardDeck{
		Title:       req.Title,
		Description: req.Description,
		Language:    language,
		Status:      status,
		CreatedBy:   utils.StringPtr(authorID),
	}

	return s.repo.CreateDeck(ctx, deck, mapFlashcardCardInputs(req.Cards))
}

func (s *FlashcardService) GetDeck(ctx context.Context, id string) (*domain.FlashcardDeckWithCards, error) {
	return s.repo.GetDeckByID(ctx, id)
}

func (s *FlashcardService) ListDecksCursor(ctx context.Context, userID string, filter domain.FlashcardDeckFilter, cursor string, limit int) ([]*domain.FlashcardDeckWithCount, *string, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	return s.repo.ListDecksCursor(ctx, userID, filter, cursor, limit)
}

func (s *FlashcardService) UpdateDeck(ctx context.Context, id string, req domain.UpdateFlashcardDeckRequest) error {
	deckInfo, err := s.repo.GetDeckByID(ctx, id)
	if err != nil {
		return err
	}

	deck := *deckInfo.FlashcardDeck
	utils.SetIfNotNil(&deck.Status, req.Status)
	if req.Title != nil {
		deck.Title = *req.Title
	}
	if req.Description != nil {
		deck.Description = req.Description
	}
	if req.Language != nil {
		deck.Language = *req.Language
	}
	deck.UpdatedBy = utils.StringPtr(ctx.Value("user_id").(string))

	return s.repo.UpdateDeck(ctx, id, deck)
}

func (s *FlashcardService) DeleteDeck(ctx context.Context, id string) error {
	return s.repo.DeleteDeck(ctx, id)
}

func (s *FlashcardService) CreateCard(ctx context.Context, deckID string, req domain.CreateFlashcardCardRequest) (*domain.FlashcardCard, error) {
	card := domain.FlashcardCard{
		DeckID:   deckID,
		Front:    req.Front,
		Back:     req.Back,
		Example:  req.Example,
		Phonetic: req.Phonetic,
	}
	return s.repo.CreateCard(ctx, card)
}

func (s *FlashcardService) UpdateCard(ctx context.Context, id string, req domain.UpdateFlashcardCardRequest) error {
	card, err := s.repo.GetCardByID(ctx, id)
	if err != nil {
		return err
	}

	updates := *card
	if req.Front != nil {
		updates.Front = *req.Front
	}
	if req.Back != nil {
		updates.Back = *req.Back
	}
	if req.Example != nil {
		updates.Example = req.Example
	}
	if req.Phonetic != nil {
		updates.Phonetic = req.Phonetic
	}

	return s.repo.UpdateCard(ctx, id, updates)
}

func (s *FlashcardService) DeleteCard(ctx context.Context, id string) error {
	return s.repo.DeleteCard(ctx, id)
}

func (s *FlashcardService) ListDueCards(ctx context.Context, userID string, deckID string, limit int) ([]*domain.FlashcardCard, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.ListDueCards(ctx, userID, deckID, limit)
}

// SubmitReview grades a single card review, applies the SM-2 scheduling
// algorithm, and persists the learner's updated per-card state.
func (s *FlashcardService) SubmitReview(ctx context.Context, userID string, cardID string, req domain.SubmitFlashcardReviewRequest) (*domain.FlashcardReviewResult, error) {
	if err := validateFlashcardRating(req.Rating); err != nil {
		return nil, err
	}

	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	review, err := s.repo.GetReview(ctx, userID, cardID)
	if err != nil {
		return nil, err
	}
	if review == nil {
		review = &domain.FlashcardReview{
			UserID:       userID,
			CardID:       cardID,
			DeckID:       card.DeckID,
			EaseFactor:   defaultEaseFactor,
			IntervalDays: 0,
			Repetitions:  0,
		}
	}

	updated := applySM2(*review, req.Rating)
	saved, err := s.repo.UpsertReview(ctx, updated)
	if err != nil {
		return nil, err
	}

	return &domain.FlashcardReviewResult{
		CardID:       saved.CardID,
		EaseFactor:   saved.EaseFactor,
		IntervalDays: saved.IntervalDays,
		Repetitions:  saved.Repetitions,
		DueAt:        saved.DueAt,
	}, nil
}

func validateFlashcardRating(rating domain.FlashcardRating) error {
	switch rating {
	case domain.FlashcardRatingAgain, domain.FlashcardRatingHard, domain.FlashcardRatingGood, domain.FlashcardRatingEasy:
		return nil
	default:
		return fmt.Errorf("invalid rating %q (expected again, hard, good, or easy)", rating)
	}
}

// applySM2 computes the next ease factor, interval, and due date for a card
// given the learner's rating, following the standard SM-2 spaced-repetition
// algorithm (Anki-style quality mapping).
func applySM2(review domain.FlashcardReview, rating domain.FlashcardRating) domain.FlashcardReview {
	now := time.Now()

	switch rating {
	case domain.FlashcardRatingAgain:
		review.Repetitions = 0
		review.IntervalDays = 1
		review.EaseFactor = math.Max(minEaseFactor, review.EaseFactor-0.20)
	case domain.FlashcardRatingHard:
		review.Repetitions++
		if review.IntervalDays <= 0 {
			review.IntervalDays = 1
		} else {
			review.IntervalDays = int(math.Round(float64(review.IntervalDays) * 1.2))
		}
		review.EaseFactor = math.Max(minEaseFactor, review.EaseFactor-0.15)
	case domain.FlashcardRatingGood:
		review.Repetitions++
		switch review.Repetitions {
		case 1:
			review.IntervalDays = 1
		case 2:
			review.IntervalDays = 6
		default:
			review.IntervalDays = int(math.Round(float64(review.IntervalDays) * review.EaseFactor))
		}
	case domain.FlashcardRatingEasy:
		review.Repetitions++
		base := review.IntervalDays
		if base <= 0 {
			base = 1
		}
		review.IntervalDays = int(math.Round(float64(base) * review.EaseFactor * 1.3))
		review.EaseFactor += 0.15
	}

	review.LastReviewedAt = &now
	r := rating
	review.LastRating = &r
	review.DueAt = now.AddDate(0, 0, review.IntervalDays)
	return review
}

// GenerateCards drafts flashcards via the configured AI provider without persisting them.
func (s *FlashcardService) GenerateCards(ctx context.Context, req domain.GenerateFlashcardsRequest) (*domain.GenerateFlashcardsResult, error) {
	if s.generator == nil {
		return nil, fmt.Errorf("AI flashcard generation is not configured (set AI_PROVIDER and AI_API_KEY)")
	}
	draft, err := s.generator.GenerateFlashcards(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(draft.Cards) == 0 {
		return nil, fmt.Errorf("ai response contained no cards")
	}
	return draft, nil
}
