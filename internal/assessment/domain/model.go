package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type Quiz struct {
	ID   uuid.UUID
	LessonID uuid.UUID
	TitleEN string
	TitleVI string
	Questions []Question
}

func (q *Quiz) Title(lang common.Language) string {
	if lang == common.LanguageVI {
		return q.TitleVI
	}
	return q.TitleEN
}

type Question struct {
	ID   uuid.UUID
	QuizID uuid.UUID
	TextEN string
	TextVI string
	Answers []Answer
}

func (q *Question) Text(lang common.Language) string {
	if lang == common.LanguageVI {
		return q.TextVI
	}
	return q.TextEN
}

type Answer struct {
	ID   uuid.UUID
	QuestionID uuid.UUID
	Text string
	IsCorrect bool
}

type QuizAttempt struct {
	ID   uuid.UUID
	QuizID uuid.UUID
	UserID uuid.UUID
	Score int
	Status string // IN_PROGRESS, COMPLETED
	CreatedAt time.Time
}

type QuizRepository interface {
	FindByLessonID(lessonID uuid.UUID, lang common.Language) (*Quiz, error)
}

type QuizAttemptRepository interface {
	Create(attempt *QuizAttempt) error
	FindByID(id uuid.UUID) (*QuizAttempt, error)
	Update(attempt *QuizAttempt) error
}
