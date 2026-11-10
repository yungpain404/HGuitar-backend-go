package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepositoryAdapter struct {
	db *sql.DB
}

func NewRefreshTokenRepositoryAdapter(db *sql.DB) *RefreshTokenRepositoryAdapter {
	return &RefreshTokenRepositoryAdapter{db: db}
}

func (r *RefreshTokenRepositoryAdapter) FindUserIDByTokenHash(tokenHash string) (uuid.UUID, error) {
	var idStr string
	err := r.db.QueryRow("SELECT user_id FROM refresh_tokens WHERE token_hash = @p1 AND expires_at > @p2 AND revoked = 0", tokenHash, time.Now()).Scan(&idStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, err
		}
		return uuid.Nil, err
	}
	return uuid.Parse(idStr)
}

func (r *RefreshTokenRepositoryAdapter) RevokeByTokenHash(tokenHash string) error {
	_, err := r.db.Exec("UPDATE refresh_tokens SET revoked = 1 WHERE token_hash = @p1", tokenHash)
	return err
}

func (r *RefreshTokenRepositoryAdapter) Save(userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked, created_at)
		VALUES (@p1, @p2, @p3, @p4, 0, @p5)`,
		uuid.New(), userID, tokenHash, expiresAt, time.Now())
	return err
}
