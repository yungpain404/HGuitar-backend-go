package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/music/domain"
)

type ChordPositionResponse struct {
	ID   string `json:"id"`
	PositionNumber int `json:"position_number"`
	BaseFret int `json:"base_fret"`
	FingersJSON string `json:"fingers_json"`
}

type ChordResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	GuitarType string `json:"guitar_type"`
	Category string `json:"category"`
	Description string `json:"description"`
	Positions []ChordPositionResponse `json:"positions"`
}

type ChordService interface {
	GetChords(guitarType, category string) ([]*ChordResponse, error)
	GetChord(id uuid.UUID) (*ChordResponse, error)
	GetPositions(chordID uuid.UUID) ([]*ChordPositionResponse, error)
}

type chordService struct {
	repo domain.ChordRepository
}

func NewChordService(repo domain.ChordRepository) ChordService {
	return &chordService{repo: repo}
}

func (s *chordService) GetChords(guitarType, category string) ([]*ChordResponse, error) {
	var chords []*domain.Chord
	var err error
	if category != "" {
		chords, err = s.repo.FindByCategory(common.ChordCategory(category))
	} else if guitarType != "" {
		chords, err = s.repo.FindByGuitarType(guitarType)
	} else {
		chords, err = s.repo.FindAll()
	}
	if err != nil {
		return nil, err
	}
	resp := make([]*ChordResponse, len(chords))
	for i, c := range chords {
		resp[i] = toChordResponse(c)
	}
	return resp, nil
}

func (s *chordService) GetChord(id uuid.UUID) (*ChordResponse, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return toChordResponse(c), nil
}

func (s *chordService) GetPositions(chordID uuid.UUID) ([]*ChordPositionResponse, error) {
	c, err := s.repo.FindByID(chordID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	resp := make([]*ChordPositionResponse, len(c.Positions))
	for i, p := range c.Positions {
		resp[i] = &ChordPositionResponse{
			ID:            p.ID.String(),
			PositionNumber: p.PositionNumber,
			BaseFret:      p.BaseFret,
			FingersJSON:   p.FingersJSON,
		}
	}
	return resp, nil
}

func toChordResponse(c *domain.Chord) *ChordResponse {
	positions := make([]ChordPositionResponse, len(c.Positions))
	for i, p := range c.Positions {
		positions[i] = ChordPositionResponse{
			ID:            p.ID.String(),
			PositionNumber: p.PositionNumber,
			BaseFret:      p.BaseFret,
			FingersJSON:   p.FingersJSON,
		}
	}
	return &ChordResponse{
		ID:          c.ID.String(),
		Name:        c.Name,
		GuitarType:  c.GuitarType,
		Category:    string(c.Category),
		Description: c.Description,
		Positions:   positions,
	}
}

type StrummingPatternResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Formula string `json:"formula"`
	Description string `json:"description"`
	Difficulty string `json:"difficulty"`
	GuitarType string `json:"guitar_type"`
	BPMMin *int `json:"bpm_min,omitempty"`
	BPMMax *int `json:"bpm_max,omitempty"`
	IsFormula bool `json:"is_formula"`
	SortOrder int `json:"sort_order"`
}

type StrummingPatternDetailResponse struct {
	StrummingPatternResponse
	Instructions []StrummingInstructionResponse `json:"instructions"`
}

type StrummingInstructionResponse struct {
	ID        string `json:"id"`
	StepNumber int   `json:"step_number"`
	Instruction string `json:"instruction"`
}

type StrummingPatternService interface {
	GetPatterns(guitarType, difficulty string) ([]*StrummingPatternResponse, error)
	GetPatternDetail(id uuid.UUID, lang string) (*StrummingPatternDetailResponse, error)
	GetPatternInstructions(id uuid.UUID, lang string) ([]*StrummingInstructionResponse, error)
	CreatePattern(req CreateStrummingPatternRequest) (*StrummingPatternResponse, error)
}

type CreateStrummingPatternRequest struct {
	Name         string `json:"name" validate:"required"`
	Formula      string `json:"formula" validate:"required"`
	Description  string `json:"description"`
	Difficulty   string `json:"difficulty" validate:"required"`
	GuitarType   string `json:"guitar_type" validate:"required"`
	BPMMin       *int   `json:"bpm_min"`
	BPMMax       *int   `json:"bpm_max"`
	IsFormula    bool   `json:"is_formula"`
	SortOrder    int    `json:"sort_order"`
}

type strummingPatternService struct {
	patternRepo domain.StrummingPatternRepository
	instrRepo   domain.StrummingInstructionRepository
}

func NewStrummingPatternService(patternRepo domain.StrummingPatternRepository, instrRepo domain.StrummingInstructionRepository) StrummingPatternService {
	return &strummingPatternService{patternRepo: patternRepo, instrRepo: instrRepo}
}

func (s *strummingPatternService) GetPatterns(guitarType, difficulty string) ([]*StrummingPatternResponse, error) {
	patterns, err := s.patternRepo.FindAll()
	if err != nil {
		return nil, err
	}
	resp := make([]*StrummingPatternResponse, 0, len(patterns))
	for _, p := range patterns {
		if guitarType != "" && p.GuitarType != guitarType {
			continue
		}
		if difficulty != "" && p.Difficulty != difficulty {
			continue
		}
		resp = append(resp, &StrummingPatternResponse{
			ID:   p.ID.String(),
			Name: p.Name,
			Formula: p.Formula,
			Description: p.Description,
			Difficulty: p.Difficulty,
			GuitarType: p.GuitarType,
			BPMMin: p.BPMMin,
			BPMMax: p.BPMMax,
			IsFormula: p.IsFormula,
			SortOrder: p.SortOrder,
		})
	}
	return resp, nil
}

func (s *strummingPatternService) GetPatternDetail(id uuid.UUID, lang string) (*StrummingPatternDetailResponse, error) {
	p, err := s.patternRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	instrs, err := s.instrRepo.FindByPatternID(id)
	if err != nil {
		return nil, err
	}
	instrResp := make([]StrummingInstructionResponse, len(instrs))
	for i, ins := range instrs {
		instrResp[i] = StrummingInstructionResponse{
			ID:        ins.ID.String(),
			StepNumber: ins.StepNumber,
			Instruction: ins.Instruction(lang),
		}
	}
	return &StrummingPatternDetailResponse{
		StrummingPatternResponse: StrummingPatternResponse{
			ID:   p.ID.String(),
			Name: p.Name,
			Formula: p.Formula,
			Description: p.Description,
			Difficulty: p.Difficulty,
			GuitarType: p.GuitarType,
			BPMMin: p.BPMMin,
			BPMMax: p.BPMMax,
			IsFormula: p.IsFormula,
			SortOrder: p.SortOrder,
		},
		Instructions: instrResp,
	}, nil
}

func (s *strummingPatternService) GetPatternInstructions(id uuid.UUID, lang string) ([]*StrummingInstructionResponse, error) {
	instrs, err := s.instrRepo.FindByPatternID(id)
	if err != nil {
		return nil, err
	}
	resp := make([]*StrummingInstructionResponse, len(instrs))
	for i, ins := range instrs {
		resp[i] = &StrummingInstructionResponse{
			ID:        ins.ID.String(),
			StepNumber: ins.StepNumber,
			Instruction: ins.Instruction(lang),
		}
	}
	return resp, nil
}

func (s *strummingPatternService) CreatePattern(req CreateStrummingPatternRequest) (*StrummingPatternResponse, error) {
	pattern := &domain.StrummingPattern{
		ID:          uuid.New(),
		Name:        req.Name,
		Formula:     req.Formula,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		GuitarType:  req.GuitarType,
		BPMMin:      req.BPMMin,
		BPMMax:      req.BPMMax,
		IsFormula:   req.IsFormula,
		Status:      common.PublicationStatusPublished,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
	}
	if err := s.patternRepo.Save(pattern); err != nil {
		return nil, err
	}
	return &StrummingPatternResponse{
		ID:   pattern.ID.String(),
		Name: pattern.Name,
		Formula: pattern.Formula,
		Description: pattern.Description,
		Difficulty: pattern.Difficulty,
		GuitarType: pattern.GuitarType,
		BPMMin: pattern.BPMMin,
		BPMMax: pattern.BPMMax,
		IsFormula: pattern.IsFormula,
		SortOrder: pattern.SortOrder,
	}, nil
}
