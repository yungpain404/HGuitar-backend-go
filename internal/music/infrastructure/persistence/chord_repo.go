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
	return r.find("")
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
		var id, name, guitarType, category, description sql.NullString
		if err := rows.Scan(&id, &name, &guitarType, &category, &description); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id.String)
		out = append(out, &domain.Chord{
			ID:          parsedID,
			Name:        name.String,
			GuitarType:  guitarType.String,
			Category:    common.ChordCategory(category.String),
			Description: description.String,
		})
	}

	// Load positions cho từng chord đã quét được.
	for _, c := range out {
		positions, err := r.loadPositions(c.ID)
		if err != nil {
			return nil, err
		}
		c.Positions = positions
	}
	return out, nil
}

// loadPositions lấy danh sách vị trí ngón tay của một hợp âm từ bảng chord_positions.
func (r *ChordRepositoryAdapter) loadPositions(chordID uuid.UUID) ([]domain.ChordPosition, error) {
	rows, err := r.db.Query(
		"SELECT id, chord_id, position_number, base_fret, fingers_json FROM chord_positions WHERE chord_id = @p1 ORDER BY position_number",
		chordID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ChordPosition
	for rows.Next() {
		var id, chordIDStr sql.NullString
		var positionNumber, baseFret sql.NullInt64
		var fingersJSON sql.NullString
		if err := rows.Scan(&id, &chordIDStr, &positionNumber, &baseFret, &fingersJSON); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id.String)
		parsedChordID, _ := uuid.Parse(chordIDStr.String)
		out = append(out, domain.ChordPosition{
			ID:           parsedID,
			ChordID:      parsedChordID,
			PositionNumber: int(positionNumber.Int64),
			BaseFret:     int(baseFret.Int64),
			FingersJSON:  fingersJSON.String,
		})
	}
	return out, nil
}
