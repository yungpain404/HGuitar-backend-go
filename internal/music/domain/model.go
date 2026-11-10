package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type Chord struct {
	ID          uuid.UUID
	Name        string
	GuitarType  string
	Category    common.ChordCategory
	Description string
	Positions   []ChordPosition
}

type ChordPosition struct {
	ID           uuid.UUID
	ChordID      uuid.UUID
	PositionNumber int
	BaseFret     int
	FingersJSON  string
}

type ChordRepository interface {
	FindAll() ([]*Chord, error)
	FindByCategory(category common.ChordCategory) ([]*Chord, error)
	FindByGuitarType(guitarType string) ([]*Chord, error)
	FindByID(id uuid.UUID) (*Chord, error)
}

type StrummingPattern struct {
	ID          uuid.UUID
	Name        string
	Formula     string
	Description string
	Difficulty  string
	GuitarType  string
	BPMMin      *int
	BPMMax      *int
	IsFormula   bool
	Status      common.PublicationStatus
	SortOrder   int
	CreatedAt   time.Time
}

type StrummingInstruction struct {
	ID             uuid.UUID
	PatternID      uuid.UUID
	StepNumber     int
	InstructionEN  string
	InstructionVI  string
}

func (i *StrummingInstruction) Instruction(lang string) string {
	if lang == "vi" {
		return i.InstructionVI
	}
	return i.InstructionEN
}

type StrummingPatternRepository interface {
	FindAll() ([]*StrummingPattern, error)
	FindByID(id uuid.UUID) (*StrummingPattern, error)
	Save(pattern *StrummingPattern) error
}

type StrummingInstructionRepository interface {
	FindByPatternID(patternID uuid.UUID) ([]*StrummingInstruction, error)
}
