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
	query := "SELECT id, title_en, title_vi, slug, description_en, description_vi, guitar_type, level, sort_order, created_at FROM courses WHERE 1=1"
	var args []any
	if guitarType != "" {
		query += " AND guitar_type = @p1"
		args = append(args, guitarType)
	}
	// build placeholder indices manually
	// go-mssqldb uses @p1..@pn; rebuild query
	return r.findCourses(query, args...)
}

func (r *CourseRepositoryAdapter) findCourses(query string, args ...any) ([]*domain.Course, error) {
	// Đơn giản hóa: go-mssqldb đón nhận @p1 theo thứ tự tham số. Ta rebuild placeholders.
	// Ở đây bỏ qua chi tiết placeholder, truyền args trực tiếp.
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Course
	for rows.Next() {
		var id string
		var c domain.Course
		var createdAt interface{}
		if err := rows.Scan(&id, &c.TitleEN, &c.TitleVI, &c.Slug, &c.DescriptionEN, &c.DescriptionVI, &c.GuitarType, &c.Level, &c.SortOrder, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		c.ID = parsedID
		out = append(out, &c)
	}
	return out, nil
}

func (r *CourseRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.Course, error) {
	list, err := r.findCourses("SELECT id, title_en, title_vi, slug, description_en, description_vi, guitar_type, level, sort_order, created_at FROM courses WHERE id = @p1", id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}
