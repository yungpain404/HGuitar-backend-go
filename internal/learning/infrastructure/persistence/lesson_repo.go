package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/learning/domain"
)

type LessonRepositoryAdapter struct {
	db *sql.DB
}

func NewLessonRepositoryAdapter(db *sql.DB) *LessonRepositoryAdapter {
	return &LessonRepositoryAdapter{db: db}
}

func (r *LessonRepositoryAdapter) FindByModuleID(moduleID uuid.UUID, lang common.Language) ([]*domain.Lesson, error) {
	return r.find("SELECT id, module_id, title_en, title_vi, content_en, content_vi, video_url, sort_order, duration_min, created_at FROM lessons WHERE module_id = @p1 ORDER BY sort_order", moduleID)
}

func (r *LessonRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.Lesson, error) {
	list, err := r.find("SELECT id, module_id, title_en, title_vi, content_en, content_vi, video_url, sort_order, duration_min, created_at FROM lessons WHERE id = @p1", id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *LessonRepositoryAdapter) FindAllByIDs(ids []uuid.UUID) ([]*domain.Lesson, error) {
	return nil, nil
}

func (r *LessonRepositoryAdapter) find(query string, args ...any) ([]*domain.Lesson, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Lesson
	for rows.Next() {
		var id, moduleID string
		var l domain.Lesson
		var createdAt interface{}
		if err := rows.Scan(&id, &moduleID, &l.TitleEN, &l.TitleVI, &l.ContentEN, &l.ContentVI, &l.VideoURL, &l.SortOrder, &l.DurationMin, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedModuleID, _ := uuid.Parse(moduleID)
		l.ID = parsedID
		l.ModuleID = parsedModuleID
		out = append(out, &l)
	}
	return out, nil
}
