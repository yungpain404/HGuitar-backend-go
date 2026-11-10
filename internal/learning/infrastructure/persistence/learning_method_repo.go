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
	return r.find("SELECT id, name_en, name_vi, description_en, description_vi, difficulty, sort_order, created_at FROM learning_methods WHERE 1=1", difficulty)
}

func (r *LearningMethodRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.LearningMethod, error) {
	list, err := r.find("SELECT id, name_en, name_vi, description_en, description_vi, difficulty, sort_order, created_at FROM learning_methods WHERE id = @p1", id)
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
		var id string
		var m domain.LearningMethod
		var createdAt interface{}
		if err := rows.Scan(&id, &m.NameEN, &m.NameVI, &m.DescriptionEN, &m.DescriptionVI, &m.Difficulty, &m.SortOrder, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		m.ID = parsedID
		out = append(out, &m)
	}
	return out, nil
}
