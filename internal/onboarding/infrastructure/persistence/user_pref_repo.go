package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/onboarding/domain"
)

type UserPreferenceRepositoryAdapter struct {
	db *sql.DB
}

func NewUserPreferenceRepositoryAdapter(db *sql.DB) *UserPreferenceRepositoryAdapter {
	return &UserPreferenceRepositoryAdapter{db: db}
}

func (r *UserPreferenceRepositoryAdapter) FindByUserID(userID uuid.UUID) (*domain.UserPreference, error) {
	row := r.db.QueryRow("SELECT id, user_id, guitar_type, completed, created_at, updated_at FROM user_preferences WHERE user_id = @p1", userID)
	var id, userIDStr, guitarType string
	var completed bool
	var createdAt, updatedAt interface{}
	if err := row.Scan(&id, &userIDStr, &guitarType, &completed, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, _ := uuid.Parse(id)
	parsedUserID, _ := uuid.Parse(userIDStr)
	return &domain.UserPreference{
		ID:   parsedID,
		UserID: parsedUserID,
		GuitarType: guitarType,
		Completed: completed,
		CreatedAt: createdAt.(time.Time),
		UpdatedAt: updatedAt.(time.Time),
	}, nil
}

func (r *UserPreferenceRepositoryAdapter) Save(pref *domain.UserPreference) error {
	_, err := r.db.Exec(`
		MERGE user_preferences AS target
		USING (SELECT @p1 AS id, @p2 AS user_id, @p3 AS guitar_type, @p4 AS completed, @p5 AS created_at, @p6 AS updated_at) AS source
		ON target.user_id = source.user_id
		WHEN MATCHED THEN UPDATE SET guitar_type = source.guitar_type, completed = source.completed, updated_at = source.updated_at
		WHEN NOT MATCHED THEN INSERT (id, user_id, guitar_type, completed, created_at, updated_at)
			VALUES (source.id, source.user_id, source.guitar_type, source.completed, source.created_at, source.updated_at);
	`, pref.ID, pref.UserID, pref.GuitarType, pref.Completed, pref.CreatedAt, pref.UpdatedAt)
	return err
}
