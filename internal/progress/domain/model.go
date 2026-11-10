package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type LessonProgress struct {
	ID   uuid.UUID
	UserID uuid.UUID
	LessonID uuid.UUID
	Status string // IN_PROGRESS, COMPLETED
	CompletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *LessonProgress) Complete() {
	p.Status = "COMPLETED"
	now := time.Now()
	p.CompletedAt = &now
	p.UpdatedAt = now
}

type LessonProgressRepository interface {
	FindByUserIDAndLessonID(userID, lessonID uuid.UUID) (*LessonProgress, error)
	FindByUserID(userID uuid.UUID) ([]*LessonProgress, error)
	Save(p *LessonProgress) error
}

type LessonRepository interface {
	FindByID(id uuid.UUID, lang common.Language) (*Lesson, error)
	FindAllByIDs(ids []uuid.UUID) ([]*Lesson, error)
}

type Lesson struct {
	ID   uuid.UUID
	TitleEN string
	TitleVI string
}

func (l *Lesson) Title(lang common.Language) string {
	if lang == common.LanguageVI {
		return l.TitleVI
	}
	return l.TitleEN
}
