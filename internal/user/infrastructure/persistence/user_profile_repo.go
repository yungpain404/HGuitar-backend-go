package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/user/domain"
)

type UserProfileRepositoryAdapter struct {
	db *sql.DB
}

func NewUserProfileRepositoryAdapter(db *sql.DB) *UserProfileRepositoryAdapter {
	return &UserProfileRepositoryAdapter{db: db}
}

func (r *UserProfileRepositoryAdapter) FindByUserID(userID uuid.UUID) (*domain.UserProfile, error) {
	row := r.db.QueryRow("SELECT id, user_id, display_name, avatar_url, guitar_type, created_at, updated_at FROM user_profiles WHERE user_id = @p1", userID)
	return scanProfile(row)
}

func (r *UserProfileRepositoryAdapter) Save(profile *domain.UserProfile) error {
	_, err := r.db.Exec(`
		MERGE user_profiles AS target
		USING (SELECT @p1 AS id, @p2 AS user_id, @p3 AS display_name, @p4 AS avatar_url, @p5 AS guitar_type, @p6 AS created_at, @p7 AS updated_at) AS source
		ON target.user_id = source.user_id
		WHEN MATCHED THEN UPDATE SET display_name = source.display_name, avatar_url = source.avatar_url, guitar_type = source.guitar_type, updated_at = source.updated_at
		WHEN NOT MATCHED THEN INSERT (id, user_id, display_name, avatar_url, guitar_type, created_at, updated_at)
			VALUES (source.id, source.user_id, source.display_name, source.avatar_url, source.guitar_type, source.created_at, source.updated_at);
	`, profile.ID, profile.UserID, profile.DisplayName, profile.AvatarURL, profile.GuitarType, profile.CreatedAt, profile.UpdatedAt)
	return err
}

func scanProfile(row *sql.Row) (*domain.UserProfile, error) {
	var id, userID string
	var displayName, avatarURL, guitarType string
	var createdAt, updatedAt time.Time
	err := row.Scan(&id, &userID, &displayName, &avatarURL, &guitarType, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, _ := uuid.Parse(id)
	parsedUserID, _ := uuid.Parse(userID)
	return domain.NewUserProfile(parsedID, parsedUserID, displayName, avatarURL, guitarType, createdAt), nil
}
