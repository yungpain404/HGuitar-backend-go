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
	return r.find(`
		SELECT
			CAST(l.id AS NVARCHAR(50)),
			CAST(l.module_id AS NVARCHAR(50)),
			l.sort_order,
			ISNULL(l.duration_minutes, 0),
			l.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.lessons l
		LEFT JOIN learning.lesson_translations ten ON ten.lesson_id = l.id AND ten.language = 'EN'
		LEFT JOIN learning.lesson_translations tvi ON tvi.lesson_id = l.id AND tvi.language = 'VI'
		WHERE l.module_id = @p1
		ORDER BY l.sort_order`, moduleID)
}

func (r *LessonRepositoryAdapter) FindByID(id uuid.UUID, lang common.Language) (*domain.Lesson, error) {
	list, err := r.find(`
		SELECT
			CAST(l.id AS NVARCHAR(50)),
			CAST(l.module_id AS NVARCHAR(50)),
			l.sort_order,
			ISNULL(l.duration_minutes, 0),
			l.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.description, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.description, '')
		FROM learning.lessons l
		LEFT JOIN learning.lesson_translations ten ON ten.lesson_id = l.id AND ten.language = 'EN'
		LEFT JOIN learning.lesson_translations tvi ON tvi.lesson_id = l.id AND tvi.language = 'VI'
		WHERE l.id = @p1`, id)
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
		var sortOrder, durationMin int
		var createdAt interface{}
		var titleEN, descEN, titleVI, descVI string
		if err := rows.Scan(&id, &moduleID, &sortOrder, &durationMin, &createdAt,
			&titleEN, &descEN, &titleVI, &descVI); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedModuleID, _ := uuid.Parse(moduleID)
		out = append(out, &domain.Lesson{
			ID:          parsedID,
			ModuleID:    parsedModuleID,
			TitleEN:     titleEN,
			TitleVI:     titleVI,
			ContentEN:   descEN,
			ContentVI:   descVI,
			SortOrder:   sortOrder,
			DurationMin: durationMin,
		})
	}
	return out, nil
}
