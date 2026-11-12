package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/learning/domain"
)

type CourseRepositoryAdapter struct {
	db *sql.DB
}

func NewCourseRepositoryAdapter(db *sql.DB) *CourseRepositoryAdapter {
	return &CourseRepositoryAdapter{db: db}
}

func (r *CourseRepositoryAdapter) FindAll(lang common.Language, guitarType string) ([]*domain.Course, error) {
	query := `
		SELECT
			CAST(c.id AS NVARCHAR(50)),
			c.guitar_type,
			c.sort_order,
			c.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.courses c
		LEFT JOIN learning.course_translations ten ON ten.course_id = c.id AND ten.language = 'EN'
		LEFT JOIN learning.course_translations tvi ON tvi.course_id = c.id AND tvi.language = 'VI'
		WHERE c.status = 'PUBLISHED'`
	var args []any
	if guitarType != "" {
		query += " AND c.guitar_type = @p1"
		args = append(args, guitarType)
	}
	query += " ORDER BY c.sort_order"
	return r.scanCourses(query, args...)
}

func (r *CourseRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.Course, error) {
	query := `
		SELECT
			CAST(c.id AS NVARCHAR(50)),
			c.guitar_type,
			c.sort_order,
			c.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.courses c
		LEFT JOIN learning.course_translations ten ON ten.course_id = c.id AND ten.language = 'EN'
		LEFT JOIN learning.course_translations tvi ON tvi.course_id = c.id AND tvi.language = 'VI'
		WHERE c.id = @p1`
	list, err := r.scanCourses(query, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *CourseRepositoryAdapter) scanCourses(query string, args ...any) ([]*domain.Course, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Course
	for rows.Next() {
		var id, guitarType string
		var sortOrder int
		var createdAt interface{}
		var titleEN, descEN, titleVI, descVI string
		if err := rows.Scan(&id, &guitarType, &sortOrder, &createdAt, &titleEN, &descEN, &titleVI, &descVI); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		out = append(out, &domain.Course{
			ID:            parsedID,
			TitleEN:       titleEN,
			TitleVI:       titleVI,
			DescriptionEN: descEN,
			DescriptionVI: descVI,
			GuitarType:    guitarType,
			SortOrder:     sortOrder,
		})
	}
	return out, nil
}
