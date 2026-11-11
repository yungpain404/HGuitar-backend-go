package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/assessment/domain"
	"github.com/hguitar/backend-go/internal/common"
)

type QuizRepositoryAdapter struct {
	db *sql.DB
}

func NewQuizRepositoryAdapter(db *sql.DB) *QuizRepositoryAdapter {
	return &QuizRepositoryAdapter{db: db}
}

func (r *QuizRepositoryAdapter) FindByLessonID(lessonID uuid.UUID, lang common.Language) (*domain.Quiz, error) {
	row := r.db.QueryRow("SELECT id, lesson_id, title_en, title_vi FROM quizzes WHERE lesson_id = @p1", lessonID)
	var id, lessonIDStr, titleEN, titleVI string
	if err := row.Scan(&id, &lessonIDStr, &titleEN, &titleVI); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, _ := uuid.Parse(id)
	return &domain.Quiz{
		ID:    parsedID,
		LessonID: parsedID,
		TitleEN: titleEN,
		TitleVI: titleVI,
	}, nil
}

type QuizAttemptRepositoryAdapter struct {
	db *sql.DB
}

func NewQuizAttemptRepositoryAdapter(db *sql.DB) *QuizAttemptRepositoryAdapter {
	return &QuizAttemptRepositoryAdapter{db: db}
}

func (r *QuizAttemptRepositoryAdapter) Create(attempt *domain.QuizAttempt) error {
	_, err := r.db.Exec(`INSERT INTO quiz_attempts (id, quiz_id, user_id, score, status, created_at) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`,
		attempt.ID, attempt.QuizID, attempt.UserID, attempt.Score, attempt.Status, attempt.CreatedAt)
	return err
}

func (r *QuizAttemptRepositoryAdapter) FindByID(id uuid.UUID) (*domain.QuizAttempt, error) {
	row := r.db.QueryRow("SELECT id, quiz_id, user_id, score, status, created_at FROM quiz_attempts WHERE id = @p1", id)
	var idStr, quizIDStr, userIDStr string
	var score int
	var status string
	var createdAt interface{}
	if err := row.Scan(&idStr, &quizIDStr, &userIDStr, &score, &status, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, _ := uuid.Parse(idStr)
	return &domain.QuizAttempt{
		ID:   parsedID,
		QuizID: parsedID,
		UserID: parsedID,
		Score: score,
		Status: status,
		CreatedAt: createdAt.(time.Time),
	}, nil
}

func (r *QuizAttemptRepositoryAdapter) Update(attempt *domain.QuizAttempt) error {
	_, err := r.db.Exec(`UPDATE quiz_attempts SET score = @p1, status = @p2 WHERE id = @p3`, attempt.Score, attempt.Status, attempt.ID)
	return err
}


