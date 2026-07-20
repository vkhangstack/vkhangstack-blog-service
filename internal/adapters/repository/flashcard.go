package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

func (u *DB) CreateDeck(ctx context.Context, deck domain.FlashcardDeck, cards []domain.FlashcardCard) (*domain.FlashcardDeck, error) {
	deck.ID = u.snowflakeNode.GenerateID()
	_, err := u.db.NewInsert().Model(&deck).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("flashcard deck not saved: %v", err)
	}

	if err := u.insertFlashcardCards(ctx, deck.ID, cards); err != nil {
		return nil, err
	}

	return &deck, nil
}

func (u *DB) GetDeckByID(ctx context.Context, id string) (*domain.FlashcardDeckWithCards, error) {
	deck := &domain.FlashcardDeck{}
	err := u.db.NewSelect().Model(deck).Where("fd.id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("flashcard deck not found: %v", err)
	}

	var cards []*domain.FlashcardCard
	err = u.db.NewSelect().
		Model(&cards).
		Where("fc.deck_id = ?", id).
		Order("fc.position ASC", "fc.id ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("flashcard cards not found: %v", err)
	}

	return &domain.FlashcardDeckWithCards{FlashcardDeck: deck, Cards: cards}, nil
}

func (u *DB) UpdateDeck(ctx context.Context, id string, updates domain.FlashcardDeck) error {
	updates.ID = id
	_, err := u.db.NewUpdate().Model(&updates).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("flashcard deck not updated: %v", err)
	}
	return nil
}

func (u *DB) DeleteDeck(ctx context.Context, id string) error {
	deck := &domain.FlashcardDeck{ID: id}
	_, err := u.db.NewDelete().Model(deck).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("flashcard deck not deleted: %v", err)
	}
	_, err = u.db.NewDelete().Model((*domain.FlashcardCard)(nil)).Where("deck_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("flashcard cards not deleted: %v", err)
	}
	return nil
}

func (u *DB) ListDecksCursor(ctx context.Context, userID string, filter domain.FlashcardDeckFilter, cursor string, limit int) ([]*domain.FlashcardDeckWithCount, *string, int, error) {
	var cursorID string
	if cursor != "" {
		id, err := utils.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, 0, err
		}
		cursorID = id
	}

	var decks []*domain.FlashcardDeck
	query := u.db.NewSelect().Model(&decks)

	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("fd.status = ?", *filter.Status)
	}
	if filter.Language != nil && *filter.Language != "" {
		query = query.Where("fd.language = ?", *filter.Language)
	}
	if filter.CreatedBy != nil && *filter.CreatedBy != "" {
		query = query.Where("fd.created_by = ?", *filter.CreatedBy)
	}
	if filter.Title != nil && *filter.Title != "" {
		query = query.Where("fd.title ILIKE ?", "%"+*filter.Title+"%")
	}

	queryCount := query.Clone()
	query = query.Order("fd.created_at DESC", "fd.id DESC")

	if cursorID != "" {
		cursorDeck := &domain.FlashcardDeck{}
		err := u.db.NewSelect().Model(cursorDeck).Where("fd.id = ?", cursorID).Limit(1).Scan(ctx)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cursor deck not found: %v", err)
		}
		query = query.Where("(fd.created_at, fd.id) < (?, ?)", cursorDeck.CreatedAt, cursorID)
	}

	err := query.Limit(limit + 1).Scan(ctx)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("flashcard decks not found: %v", err)
	}

	var nextCursor *string
	if len(decks) > limit {
		decks = decks[:limit]
		nextCursor = utils.StringPtr(utils.EncodeCursor(decks[len(decks)-1].ID))
	}
	total, _ := queryCount.Count(ctx)

	var deckIDs []string
	for _, deck := range decks {
		deckIDs = append(deckIDs, deck.ID)
	}

	cardCounts, err := u.getFlashcardCounts(ctx, deckIDs)
	if err != nil {
		return nil, nil, 0, err
	}
	dueCounts, err := u.getFlashcardDueCounts(ctx, userID, deckIDs)
	if err != nil {
		return nil, nil, 0, err
	}

	result := []*domain.FlashcardDeckWithCount{}
	for _, deck := range decks {
		result = append(result, &domain.FlashcardDeckWithCount{
			FlashcardDeck: deck,
			CardCount:     cardCounts[deck.ID],
			DueCount:      dueCounts[deck.ID],
		})
	}
	return result, nextCursor, total, nil
}

func (u *DB) CreateCard(ctx context.Context, card domain.FlashcardCard) (*domain.FlashcardCard, error) {
	card.ID = u.snowflakeNode.GenerateID()

	var maxPos int
	err := u.db.NewSelect().
		Model((*domain.FlashcardCard)(nil)).
		ColumnExpr("COALESCE(MAX(position), -1)").
		Where("deck_id = ?", card.DeckID).
		Scan(ctx, &maxPos)
	if err != nil {
		return nil, fmt.Errorf("failed to compute card position: %v", err)
	}
	card.Position = maxPos + 1

	_, err = u.db.NewInsert().Model(&card).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("flashcard card not saved: %v", err)
	}
	return &card, nil
}

func (u *DB) GetCardByID(ctx context.Context, id string) (*domain.FlashcardCard, error) {
	card := &domain.FlashcardCard{}
	err := u.db.NewSelect().Model(card).Where("fc.id = ?", id).Limit(1).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("flashcard card not found: %v", err)
	}
	return card, nil
}

func (u *DB) UpdateCard(ctx context.Context, id string, updates domain.FlashcardCard) error {
	updates.ID = id
	_, err := u.db.NewUpdate().Model(&updates).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("flashcard card not updated: %v", err)
	}
	return nil
}

func (u *DB) DeleteCard(ctx context.Context, id string) error {
	card := &domain.FlashcardCard{ID: id}
	_, err := u.db.NewDelete().Model(card).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("flashcard card not deleted: %v", err)
	}
	return nil
}

func (u *DB) ListDueCards(ctx context.Context, userID string, deckID string, limit int) ([]*domain.FlashcardCard, error) {
	var cards []*domain.FlashcardCard
	err := u.db.NewSelect().
		Model(&cards).
		ColumnExpr("fc.*").
		Join("LEFT JOIN flashcard_reviews AS fr ON fr.card_id = fc.id AND fr.user_id = ?", userID).
		Where("fc.deck_id = ?", deckID).
		Where("(fr.due_at IS NULL OR fr.due_at <= NOW())").
		OrderExpr("fr.due_at ASC NULLS FIRST").
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("due flashcards not found: %v", err)
	}
	return cards, nil
}

func (u *DB) GetReview(ctx context.Context, userID string, cardID string) (*domain.FlashcardReview, error) {
	review := &domain.FlashcardReview{}
	err := u.db.NewSelect().
		Model(review).
		Where("fr.user_id = ? AND fr.card_id = ?", userID, cardID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("flashcard review not found: %v", err)
	}
	return review, nil
}

func (u *DB) UpsertReview(ctx context.Context, review domain.FlashcardReview) (*domain.FlashcardReview, error) {
	if review.ID == "" {
		review.ID = u.snowflakeNode.GenerateID()
	}
	_, err := u.db.NewInsert().
		Model(&review).
		On("CONFLICT (user_id, card_id) DO UPDATE").
		Set("ease_factor = EXCLUDED.ease_factor").
		Set("interval_days = EXCLUDED.interval_days").
		Set("repetitions = EXCLUDED.repetitions").
		Set("due_at = EXCLUDED.due_at").
		Set("last_reviewed_at = EXCLUDED.last_reviewed_at").
		Set("last_rating = EXCLUDED.last_rating").
		Set("updated_at = now()").
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("flashcard review not saved: %v", err)
	}
	return &review, nil
}

func (u *DB) insertFlashcardCards(ctx context.Context, deckID string, cards []domain.FlashcardCard) error {
	if len(cards) == 0 {
		return nil
	}
	for i := range cards {
		cards[i].ID = u.snowflakeNode.GenerateID()
		cards[i].DeckID = deckID
		cards[i].Position = i
	}
	_, err := u.db.NewInsert().Model(&cards).Exec(ctx)
	if err != nil {
		return fmt.Errorf("flashcard cards not saved: %v", err)
	}
	return nil
}

func (u *DB) getFlashcardCounts(ctx context.Context, deckIDs []string) (map[string]int, error) {
	if len(deckIDs) == 0 {
		return map[string]int{}, nil
	}

	var rows []struct {
		DeckID string `bun:"deck_id"`
		Count  int    `bun:"count"`
	}

	err := u.db.NewSelect().
		Model((*domain.FlashcardCard)(nil)).
		ColumnExpr("deck_id, COUNT(*) AS count").
		Where("deck_id IN (?)", bun.In(deckIDs)).
		Group("deck_id").
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get flashcard counts: %w", err)
	}

	counts := make(map[string]int, len(deckIDs))
	for _, r := range rows {
		counts[r.DeckID] = r.Count
	}
	return counts, nil
}

func (u *DB) getFlashcardDueCounts(ctx context.Context, userID string, deckIDs []string) (map[string]int, error) {
	if len(deckIDs) == 0 {
		return map[string]int{}, nil
	}

	var rows []struct {
		DeckID string `bun:"deck_id"`
		Count  int    `bun:"count"`
	}

	err := u.db.NewSelect().
		TableExpr("flashcard_cards AS fc").
		Join("LEFT JOIN flashcard_reviews AS fr ON fr.card_id = fc.id AND fr.user_id = ?", userID).
		ColumnExpr("fc.deck_id AS deck_id, COUNT(*) AS count").
		Where("fc.deck_id IN (?)", bun.In(deckIDs)).
		Where("(fr.due_at IS NULL OR fr.due_at <= NOW())").
		Group("fc.deck_id").
		Scan(ctx, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to get flashcard due counts: %w", err)
	}

	counts := make(map[string]int, len(deckIDs))
	for _, r := range rows {
		counts[r.DeckID] = r.Count
	}
	return counts, nil
}
