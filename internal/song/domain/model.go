package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type Song struct {
	ID   uuid.UUID
	TitleEN string
	TitleVI string
	Artist string
	Source string
	Difficulty string
	Language string
	CreatedAt time.Time
}

func (s *Song) Title(lang common.Language) string {
	if lang == common.LanguageVI {
		return s.TitleVI
	}
	return s.TitleEN
}

type SongSection struct {
	ID   uuid.UUID
	SongID uuid.UUID
	Name string
	Chords string
	Order int
}

type SongChord struct {
	ID   uuid.UUID
	SongID uuid.UUID
	ChordID uuid.UUID
}

type SongRepository interface {
	FindAll(source, difficulty, lang string) ([]*Song, error)
	FindByID(id uuid.UUID, lang string) (*Song, error)
	Save(song *Song) error
}
