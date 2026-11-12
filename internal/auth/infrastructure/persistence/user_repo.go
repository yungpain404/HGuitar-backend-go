package persistence

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/auth/domain"
)

type UserRepositoryAdapter struct {
	db *sql.DB
}

func NewUserRepositoryAdapter(db *sql.DB) *UserRepositoryAdapter {
	return &UserRepositoryAdapter{db: db}
}

func (r *UserRepositoryAdapter) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT CASE WHEN COUNT(*) > 0 THEN 1 ELSE 0 END FROM auth.users WHERE email = @p1", email).Scan(&exists)
	return exists, err
}

func (r *UserRepositoryAdapter) FindByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT CAST(id AS NVARCHAR(50)), email, password_hash, is_active, created_at, updated_at FROM auth.users WHERE email = @p1", email)
	return scanUser(row)
}

func (r *UserRepositoryAdapter) FindByID(id uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRow("SELECT CAST(id AS NVARCHAR(50)), email, password_hash, is_active, created_at, updated_at FROM auth.users WHERE id = @p1", id)
	return scanUser(row)
}

func (r *UserRepositoryAdapter) Save(user *domain.User) error {
	_, err := r.db.Exec(`
		MERGE auth.users AS target
		USING (SELECT @p1 AS id, @p2 AS email, @p3 AS password_hash, @p4 AS is_active, @p5 AS created_at, @p6 AS updated_at) AS source
		ON target.id = source.id
		WHEN MATCHED THEN UPDATE SET email = source.email, password_hash = source.password_hash, is_active = source.is_active, updated_at = source.updated_at
		WHEN NOT MATCHED THEN INSERT (id, email, password_hash, is_active, created_at, updated_at)
			VALUES (source.id, source.email, source.password_hash, source.is_active, source.created_at, source.updated_at);
	`, user.ID, user.Email, user.PasswordHash, user.Active, user.CreatedAt, user.UpdatedAt)
	return err
}

func scanUser(row *sql.Row) (*domain.User, error) {
	var id, email, passwordHash string
	var active bool
	var createdAt, updatedAt time.Time
	err := row.Scan(&id, &email, &passwordHash, &active, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return domain.NewUser(parsedID, email, passwordHash, active, nil, createdAt), nil
}
