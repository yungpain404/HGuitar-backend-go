package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/progress/domain"
)

type LessonProgressRepositoryAdapter struct {
	db *sql.DB
}

func NewLessonProgressRepositoryAdapter(db *sql.DB) *LessonProgressRepositoryAdapter {
	return &LessonProgressRepositoryAdapter{db: db}
}

func (r *LessonProgressRepositoryAdapter) FindByUserIDAndLessonID(userID, lessonID uuid.UUID) (*domain.LessonProgress, error) {
	row := r.db.QueryRow("SELECT id, user_id, lesson_id, status, completed_at, created_at, updated_at FROM lesson_progress WHERE user_id = @p1 AND lesson_id = @p2", userID, lessonID)
	return scanProgress(row)
}

func (r *LessonProgressRepositoryAdapter) FindByUserID(userID uuid.UUID) ([]*domain.LessonProgress, error) {
	rows, err := r.db.Query("SELECT id, user_id, lesson_id, status, completed_at, created_at, updated_at FROM lesson_progress WHERE user_id = @p1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.LessonProgress
	for rows.Next() {
		p, err := scanProgressRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *LessonProgressRepositoryAdapter) Save(p *domain.LessonProgress) error {
	_, err := r.db.Exec(`
		MERGE lesson_progress AS target
		USING (SELECT @p1 AS id, @p2 AS user_id, @p3 AS lesson_id, @p4 AS status, @p5 AS completed_at, @p6 AS created_at, @p7 AS updated_at) AS source
		ON target.id = source.id
		WHEN MATCHED THEN UPDATE SET status = source.status, completed_at = source.completed_at, updated_at = source.updated_at
		WHEN NOT MATCHED THEN INSERT (id, user_id, lesson_id, status, completed_at, created_at, updated_at)
			VALUES (source.id, source.user_id, source.lesson_id, source.status, source.completed_at, source.created_at, source.updated_at);
	`, p.ID, p.UserID, p.LessonID, p.Status, p.CompletedAt, p.CreatedAt, p.UpdatedAt)
	return err
}

func scanProgress(row *sql.Row) (*domain.LessonProgress, error) {
	var id, userID, lessonID string
	var status string
	var completedAt, createdAt, updatedAt interface{}
	err := row.Scan(&id, &userID, &lessonID, &status, &completedAt, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, _ := uuid.Parse(id)
	parsedUserID, _ := uuid.Parse(userID)
	parsedLessonID, _ := uuid.Parse(lessonID)
	var completedAtTime *time.Time
	if completedAt != nil {
		if t, ok := completedAt.(time.Time); ok {
			completedAtTime = &t
		}
	}
	return &domain.LessonProgress{
		ID:   parsedID,
		UserID: parsedUserID,
		LessonID: parsedLessonID,
		Status: status,
		CompletedAt: completedAtTime,
		CreatedAt: createdAt.(time.Time),
		UpdatedAt: updatedAt.(time.Time),
	}, nil
}

func scanProgressRow(rows *sql.Rows) (*domain.LessonProgress, error) {
	var id, userID, lessonID string
	var status string
	var completedAt, createdAt, updatedAt interface{}
	if err := rows.Scan(&id, &userID, &lessonID, &status, &completedAt, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	parsedID, _ := uuid.Parse(id)
	parsedUserID, _ := uuid.Parse(userID)
	parsedLessonID, _ := uuid.Parse(lessonID)
	var completedAtTime *time.Time
	if completedAt != nil {
		if t, ok := completedAt.(time.Time); ok {
			completedAtTime = &t
		}
	}
	return &domain.LessonProgress{
		ID:   parsedID,
		UserID: parsedUserID,
		LessonID: parsedLessonID,
		Status: status,
		CompletedAt: completedAtTime,
		CreatedAt: createdAt.(time.Time),
		UpdatedAt: updatedAt.(time.Time),
	}, nil
}
