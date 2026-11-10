package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type Course struct {
	ID   uuid.UUID
	TitleEN string
	TitleVI string
	Slug  string
	DescriptionEN string
	DescriptionVI string
	GuitarType string
	Level string
	SortOrder int
	CreatedAt time.Time
}

func (c *Course) Title(lang common.Language) string {
	if lang == common.LanguageVI {
		return c.TitleVI
	}
	return c.TitleEN
}

func (c *Course) Description(lang common.Language) string {
	if lang == common.LanguageVI {
		return c.DescriptionVI
	}
	return c.DescriptionEN
}

type Module struct {
	ID   uuid.UUID
	CourseID uuid.UUID
	TitleEN string
	TitleVI string
	DescriptionEN string
	DescriptionVI string
	SortOrder int
	CreatedAt time.Time
}

func (m *Module) Title(lang common.Language) string {
	if lang == common.LanguageVI {
		return m.TitleVI
	}
	return m.TitleEN
}

type Lesson struct {
	ID   uuid.UUID
	ModuleID uuid.UUID
	TitleEN string
	TitleVI string
	ContentEN string
	ContentVI string
	VideoURL string
	SortOrder int
	DurationMin int
	CreatedAt time.Time
}

func (l *Lesson) Title(lang common.Language) string {
	if lang == common.LanguageVI {
		return l.TitleVI
	}
	return l.TitleEN
}

type LearningMethod struct {
	ID   uuid.UUID
	NameEN string
	NameVI string
	DescriptionEN string
	DescriptionVI string
	Difficulty string
	SortOrder int
	CreatedAt time.Time
}

func (lm *LearningMethod) Name(lang common.Language) string {
	if lang == common.LanguageVI {
		return lm.NameVI
	}
	return lm.NameEN
}

type CourseRepository interface {
	FindAll(lang common.Language, guitarType string) ([]*Course, error)
	FindByID(id uuid.UUID, lang common.Language) (*Course, error)
}

type ModuleRepository interface {
	FindByID(id uuid.UUID, lang common.Language) (*Module, error)
	FindByCourseID(courseID uuid.UUID, lang common.Language) ([]*Module, error)
}

type LessonRepository interface {
	FindByModuleID(moduleID uuid.UUID, lang common.Language) ([]*Lesson, error)
	FindByID(id uuid.UUID, lang common.Language) (*Lesson, error)
	FindAllByIDs(ids []uuid.UUID) ([]*Lesson, error)
}

type LearningMethodRepository interface {
	FindAll(difficulty string, lang common.Language) ([]*LearningMethod, error)
	FindByID(id uuid.UUID, lang common.Language) (*LearningMethod, error)
}
