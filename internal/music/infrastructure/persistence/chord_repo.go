package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/music/domain"
)

type ChordRepositoryAdapter struct {
	db *sql.DB
}

func NewChordRepositoryAdapter(db *sql.DB) *ChordRepositoryAdapter {
	return &ChordRepositoryAdapter{db: db}
}

func (r *ChordRepositoryAdapter) FindAll() ([]*domain.Chord, error) {
	rows, err := r.db.Query("SELECT id, name, guitar_type, category, description FROM chords")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Chord
	for rows.Next() {
		var id, name, guitarType, category, description string
		if err := rows.Scan(&id, &name, &guitarType, &category, &description); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		out = append(out, &domain.Chord{
			ID:          parsedID,
			Name:        name,
			GuitarType:  guitarType,
			Category:    common.ChordCategory(category),
			Description: description,
		})
	}
	return out, nil
}

func (r *ChordRepositoryAdapter) FindByCategory(category common.ChordCategory) ([]*domain.Chord, error) {
	return r.find("WHERE category = @p1", category)
}

func (r *ChordRepositoryAdapter) FindByGuitarType(guitarType string) ([]*domain.Chord, error) {
	return r.find("WHERE guitar_type = @p1", guitarType)
}

func (r *ChordRepositoryAdapter) FindByID(id uuid.UUID) (*domain.Chord, error) {
	rows, err := r.find("WHERE id = @p1", id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}

func (r *ChordRepositoryAdapter) find(where string, args ...any) ([]*domain.Chord, error) {
	query := "SELECT id, name, guitar_type, category, description FROM chords " + where
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Chord
	for rows.Next() {
		var id, name, guitarType, category, description string
		if err := rows.Scan(&id, &name, &guitarType, &category, &description); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		out = append(out, &domain.Chord{
			ID:          parsedID,
			Name:        name,
			GuitarType:  guitarType,
			Category:    common.ChordCategory(category),
			Description: description,
		})
	}
	return out, nil
}
