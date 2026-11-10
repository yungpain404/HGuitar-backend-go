package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/learning/domain"
)

type ModuleRepositoryAdapter struct {
	db *sql.DB
}

func NewModuleRepositoryAdapter(db *sql.DB) *ModuleRepositoryAdapter {
	return &ModuleRepositoryAdapter{db: db}
}

func (r *ModuleRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.Module, error) {
	list, err := r.find("SELECT id, course_id, title_en, title_vi, description_en, description_vi, sort_order, created_at FROM modules WHERE id = @p1", id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *ModuleRepositoryAdapter) FindByCourseID(courseID uuid.UUID, lang common.Language) ([]*domain.Module, error) {
	return r.find("SELECT id, course_id, title_en, title_vi, description_en, description_vi, sort_order, created_at FROM modules WHERE course_id = @p1 ORDER BY sort_order", courseID)
}

func (r *ModuleRepositoryAdapter) find(query string, args ...any) ([]*domain.Module, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Module
	for rows.Next() {
		var id, courseID string
		var m domain.Module
		var createdAt interface{}
		if err := rows.Scan(&id, &courseID, &m.TitleEN, &m.TitleVI, &m.DescriptionEN, &m.DescriptionVI, &m.SortOrder, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedCourseID, _ := uuid.Parse(courseID)
		m.ID = parsedID
		m.CourseID = parsedCourseID
		out = append(out, &m)
	}
	return out, nil
}
