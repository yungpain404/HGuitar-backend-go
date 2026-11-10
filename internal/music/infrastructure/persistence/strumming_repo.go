package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/music/domain"
)

type StrummingRepositoryAdapter struct {
	db *sql.DB
}

func NewStrummingRepositoryAdapter(db *sql.DB) *StrummingRepositoryAdapter {
	return &StrummingRepositoryAdapter{db: db}
}

func (r *StrummingRepositoryAdapter) FindAll() ([]*domain.StrummingPattern, error) {
	rows, err := r.db.Query("SELECT id, name, formula, description, difficulty, guitar_type, bpm_min, bpm_max, is_formula, status, sort_order, created_at FROM strumming_patterns")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StrummingPattern
	for rows.Next() {
		var id, name, formula, description, difficulty, guitarType, status string
		var bpmMin, bpmMax *int
		var isFormula bool
		var sortOrder int
		var createdAt interface{}
		if err := rows.Scan(&id, &name, &formula, &description, &difficulty, &guitarType, &bpmMin, &bpmMax, &isFormula, &status, &sortOrder, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		out = append(out, &domain.StrummingPattern{
			ID:          parsedID,
			Name:        name,
			Formula:     formula,
			Description: description,
			Difficulty:  difficulty,
			GuitarType:  guitarType,
			BPMMin:      bpmMin,
			BPMMax:      bpmMax,
			IsFormula:   isFormula,
			Status:      common.PublicationStatus(status),
			SortOrder:   sortOrder,
		})
	}
	return out, nil
}

func (r *StrummingRepositoryAdapter) FindByID(id uuid.UUID) (*domain.StrummingPattern, error) {
	rows, err := r.db.Query("SELECT id, name, formula, description, difficulty, guitar_type, bpm_min, bpm_max, is_formula, status, sort_order, created_at FROM strumming_patterns WHERE id = @p1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var idStr, name, formula, description, difficulty, guitarType, status string
		var bpmMin, bpmMax *int
		var isFormula bool
		var sortOrder int
		var createdAt interface{}
		if err := rows.Scan(&idStr, &name, &formula, &description, &difficulty, &guitarType, &bpmMin, &bpmMax, &isFormula, &status, &sortOrder, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(idStr)
		return &domain.StrummingPattern{
			ID:          parsedID,
			Name:        name,
			Formula:     formula,
			Description: description,
			Difficulty:  difficulty,
			GuitarType:  guitarType,
			BPMMin:      bpmMin,
			BPMMax:      bpmMax,
			IsFormula:   isFormula,
			Status:      common.PublicationStatus(status),
			SortOrder:   sortOrder,
		}, nil
	}
	return nil, nil
}

func (r *StrummingRepositoryAdapter) Save(pattern *domain.StrummingPattern) error {
	_, err := r.db.Exec(`
		INSERT INTO strumming_patterns (id, name, formula, description, difficulty, guitar_type, bpm_min, bpm_max, is_formula, status, sort_order, created_at)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12)`,
		pattern.ID, pattern.Name, pattern.Formula, pattern.Description, pattern.Difficulty,
		pattern.GuitarType, pattern.BPMMin, pattern.BPMMax, pattern.IsFormula, pattern.Status, pattern.SortOrder, pattern.CreatedAt)
	return err
}

type StrummingInstructionRepositoryAdapter struct {
	db *sql.DB
}

func NewStrummingInstructionRepositoryAdapter(db *sql.DB) *StrummingInstructionRepositoryAdapter {
	return &StrummingInstructionRepositoryAdapter{db: db}
}

func (r *StrummingInstructionRepositoryAdapter) FindByPatternID(patternID uuid.UUID) ([]*domain.StrummingInstruction, error) {
	rows, err := r.db.Query("SELECT id, pattern_id, step_number, instruction_en, instruction_vi FROM strumming_instructions WHERE pattern_id = @p1 ORDER BY step_number", patternID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.StrummingInstruction
	for rows.Next() {
		var id, patternID string
		var stepNumber int
		var instrEN, instrVI string
		if err := rows.Scan(&id, &patternID, &stepNumber, &instrEN, &instrVI); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedPatternID, _ := uuid.Parse(patternID)
		out = append(out, &domain.StrummingInstruction{
			ID:            parsedID,
			PatternID:     parsedPatternID,
			StepNumber:    stepNumber,
			InstructionEN: instrEN,
			InstructionVI: instrVI,
		})
	}
	return out, nil
}
