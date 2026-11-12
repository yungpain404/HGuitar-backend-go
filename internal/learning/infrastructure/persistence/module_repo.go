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
	list, err := r.find(`
		SELECT
			CAST(m.id AS NVARCHAR(50)),
			CAST(m.course_id AS NVARCHAR(50)),
			m.sort_order,
			m.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.modules m
		LEFT JOIN learning.module_translations ten ON ten.module_id = m.id AND ten.language = 'EN'
		LEFT JOIN learning.module_translations tvi ON tvi.module_id = m.id AND tvi.language = 'VI'
		WHERE m.id = @p1`, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *ModuleRepositoryAdapter) FindByCourseID(courseID uuid.UUID, lang common.Language) ([]*domain.Module, error) {
	return r.find(`
		SELECT
			CAST(m.id AS NVARCHAR(50)),
			CAST(m.course_id AS NVARCHAR(50)),
			m.sort_order,
			m.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.modules m
		LEFT JOIN learning.module_translations ten ON ten.module_id = m.id AND ten.language = 'EN'
		LEFT JOIN learning.module_translations tvi ON tvi.module_id = m.id AND tvi.language = 'VI'
		WHERE m.course_id = @p1
		ORDER BY m.sort_order`, courseID)
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
		var sortOrder int
		var createdAt interface{}
		var titleEN, descEN, titleVI, descVI string
		if err := rows.Scan(&id, &courseID, &sortOrder, &createdAt, &titleEN, &descEN, &titleVI, &descVI); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedCourseID, _ := uuid.Parse(courseID)
		out = append(out, &domain.Module{
			ID:            parsedID,
			CourseID:      parsedCourseID,
			TitleEN:       titleEN,
			TitleVI:       titleVI,
			DescriptionEN: descEN,
			DescriptionVI: descVI,
			SortOrder:     sortOrder,
		})
	}
	return out, nil
}
