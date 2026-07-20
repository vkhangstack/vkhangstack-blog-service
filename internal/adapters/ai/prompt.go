package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

// aiQuizPayload is the JSON shape every provider is instructed to return.
type aiQuizPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Questions   []struct {
		Prompt       string   `json:"prompt"`
		Options      []string `json:"options"`
		CorrectIndex int      `json:"correct_index"`
		Explanation  string   `json:"explanation"`
	} `json:"questions"`
}

func buildPrompt(req domain.GenerateQuizRequest) (system string, user string) {
	count := req.Count
	if count <= 0 {
		count = 10
	}
	if count > 30 {
		count = 30
	}
	examType := req.ExamType
	if examType == "" {
		examType = domain.QuizExamTypeGeneral
	}

	system = `You are an expert English test-prep author specializing in IELTS and TOEIC exam preparation. You write realistic, exam-style multiple-choice questions with plausible distractors and concise explanations.

Respond with a single JSON object only — no markdown fences, no commentary. Schema:
{
  "title": string,            // short quiz title
  "description": string,      // 1-2 sentence description
  "questions": [
    {
      "prompt": string,        // the question (include the sentence with a blank ___ where relevant)
      "options": [string],     // exactly 4 answer options, plain text without A/B/C/D prefixes
      "correct_index": number, // 0-based index of the correct option
      "explanation": string    // brief explanation of why the answer is correct
    }
  ]
}`

	var sb strings.Builder
	fmt.Fprintf(&sb, "Generate %d multiple-choice questions about: %s.\n", count, req.Topic)
	fmt.Fprintf(&sb, "Exam style: %s.\n", strings.ToUpper(string(examType)))
	if req.Skill != nil && *req.Skill != "" {
		fmt.Fprintf(&sb, "Skill focus: %s.\n", *req.Skill)
	}
	if req.Level != nil && *req.Level != "" {
		fmt.Fprintf(&sb, "Difficulty level: %s.\n", *req.Level)
	}
	sb.WriteString("Vary the question forms and make distractors realistic.")
	user = sb.String()
	return system, user
}

// parseQuizPayload extracts and validates the JSON object from a model response.
func parseQuizPayload(raw string, req domain.GenerateQuizRequest) (*domain.CreateQuizRequest, error) {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end <= start {
		return nil, fmt.Errorf("ai response contains no JSON object")
	}

	var payload aiQuizPayload
	if err := json.Unmarshal([]byte(raw[start:end+1]), &payload); err != nil {
		return nil, fmt.Errorf("ai response is not valid JSON: %v", err)
	}
	if payload.Title == "" || len(payload.Questions) == 0 {
		return nil, fmt.Errorf("ai response is missing title or questions")
	}

	examType := req.ExamType
	if examType == "" {
		examType = domain.QuizExamTypeGeneral
	}
	status := domain.QuizStatusDraft

	result := &domain.CreateQuizRequest{
		Title:    payload.Title,
		ExamType: examType,
		Skill:    req.Skill,
		Status:   status,
	}
	if payload.Description != "" {
		result.Description = &payload.Description
	}

	for i, q := range payload.Questions {
		if q.Prompt == "" || len(q.Options) < 2 {
			return nil, fmt.Errorf("ai question %d is malformed", i+1)
		}
		if q.CorrectIndex < 0 || q.CorrectIndex >= len(q.Options) {
			return nil, fmt.Errorf("ai question %d has correct_index out of range", i+1)
		}
		input := domain.QuizQuestionInput{
			Prompt:       q.Prompt,
			Options:      q.Options,
			CorrectIndex: q.CorrectIndex,
		}
		if q.Explanation != "" {
			explanation := q.Explanation
			input.Explanation = &explanation
		}
		result.Questions = append(result.Questions, input)
	}
	return result, nil
}

// aiFlashcardPayload is the JSON shape every provider is instructed to return.
type aiFlashcardPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Cards       []struct {
		Front    string `json:"front"`
		Back     string `json:"back"`
		Example  string `json:"example"`
		Phonetic string `json:"phonetic"`
	} `json:"cards"`
}

func buildFlashcardPrompt(req domain.GenerateFlashcardsRequest) (system string, user string) {
	count := req.Count
	if count <= 0 {
		count = 10
	}
	if count > 30 {
		count = 30
	}
	language := req.Language
	if language == "" {
		language = "en"
	}

	system = `You are an expert English vocabulary teacher who writes concise, memorable flashcards for language learners.

Respond with a single JSON object only — no markdown fences, no commentary. Schema:
{
  "title": string,            // short deck title
  "description": string,      // 1-2 sentence description
  "cards": [
    {
      "front": string,        // the English word or phrase to learn
      "back": string,         // the meaning/translation/definition
      "example": string,      // a short example sentence using the word (optional but preferred)
      "phonetic": string      // IPA or simplified pronunciation (optional)
    }
  ]
}`

	var sb strings.Builder
	fmt.Fprintf(&sb, "Generate %d English flashcards about: %s.\n", count, req.Topic)
	fmt.Fprintf(&sb, "Target language for definitions: %s.\n", language)
	if req.Level != nil && *req.Level != "" {
		fmt.Fprintf(&sb, "Difficulty level: %s.\n", *req.Level)
	}
	sb.WriteString("Keep fronts short (a word or short phrase) and backs concise. Avoid duplicates.")
	user = sb.String()
	return system, user
}

// parseFlashcardPayload extracts and validates the JSON object from a model response.
func parseFlashcardPayload(raw string) (*domain.GenerateFlashcardsResult, error) {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start < 0 || end <= start {
		return nil, fmt.Errorf("ai response contains no JSON object")
	}

	var payload aiFlashcardPayload
	if err := json.Unmarshal([]byte(raw[start:end+1]), &payload); err != nil {
		return nil, fmt.Errorf("ai response is not valid JSON: %v", err)
	}
	if payload.Title == "" || len(payload.Cards) == 0 {
		return nil, fmt.Errorf("ai response is missing title or cards")
	}

	result := &domain.GenerateFlashcardsResult{
		Title:       payload.Title,
		Description: payload.Description,
	}

	for i, c := range payload.Cards {
		if c.Front == "" || c.Back == "" {
			return nil, fmt.Errorf("ai card %d is malformed", i+1)
		}
		input := domain.FlashcardCardInput{
			Front: c.Front,
			Back:  c.Back,
		}
		if c.Example != "" {
			example := c.Example
			input.Example = &example
		}
		if c.Phonetic != "" {
			phonetic := c.Phonetic
			input.Phonetic = &phonetic
		}
		result.Cards = append(result.Cards, input)
	}
	return result, nil
}
