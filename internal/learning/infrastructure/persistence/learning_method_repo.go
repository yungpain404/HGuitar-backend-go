package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/learning/domain"
)

type LearningMethodRepositoryAdapter struct {
	db *sql.DB
}

func NewLearningMethodRepositoryAdapter(db *sql.DB) *LearningMethodRepositoryAdapter {
	return &LearningMethodRepositoryAdapter{db: db}
}

func (r *LearningMethodRepositoryAdapter) FindAll(difficulty string, lang common.Language) ([]*domain.LearningMethod, error) {
	query := `
		SELECT
			CAST(m.id AS NVARCHAR(50)),
			m.difficulty,
			m.sort_order,
			m.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.learning_methods m
		LEFT JOIN learning.learning_method_translations ten ON ten.method_id = m.id AND ten.language = 'EN'
		LEFT JOIN learning.learning_method_translations tvi ON tvi.method_id = m.id AND tvi.language = 'VI'
		WHERE m.status = 'PUBLISHED'`
	var args []any
	if difficulty != "" {
		query += " AND m.difficulty = @p1"
		args = append(args, difficulty)
	}
	query += " ORDER BY m.sort_order"
	return r.find(query, args...)
}

func (r *LearningMethodRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.LearningMethod, error) {
	list, err := r.find(`
		SELECT
			CAST(m.id AS NVARCHAR(50)),
			m.difficulty,
			m.sort_order,
			m.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.learning_methods m
		LEFT JOIN learning.learning_method_translations ten ON ten.method_id = m.id AND ten.language = 'EN'
		LEFT JOIN learning.learning_method_translations tvi ON tvi.method_id = m.id AND tvi.language = 'VI'
		WHERE m.id = @p1`, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *LearningMethodRepositoryAdapter) find(query string, args ...any) ([]*domain.LearningMethod, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.LearningMethod
	for rows.Next() {
		var id, difficulty string
		var sortOrder int
		var createdAt interface{}
		var nameEN, descEN, nameVI, descVI string
		if err := rows.Scan(&id, &difficulty, &sortOrder, &createdAt,
			&nameEN, &descEN, &nameVI, &descVI); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		out = append(out, &domain.LearningMethod{
			ID:            parsedID,
			NameEN:        nameEN,
			NameVI:        nameVI,
			DescriptionEN: descEN,
			DescriptionVI: descVI,
			Difficulty:    difficulty,
			SortOrder:     sortOrder,
		})
	}
	return out, nil
}
