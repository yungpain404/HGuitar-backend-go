package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/song/domain"
)

type SongRepositoryAdapter struct {
	db *sql.DB
}

func NewSongRepositoryAdapter(db *sql.DB) *SongRepositoryAdapter {
	return &SongRepositoryAdapter{db: db}
}

func (r *SongRepositoryAdapter) FindAll(source, difficulty, lang string) ([]*domain.Song, error) {
	rows, err := r.db.Query("SELECT id, title_en, title_vi, artist, source, difficulty, language, created_at FROM songs WHERE 1=1", source, difficulty)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Song
	for rows.Next() {
		var id string
		var s domain.Song
		var createdAt interface{}
		if err := rows.Scan(&id, &s.TitleEN, &s.TitleVI, &s.Artist, &s.Source, &s.Difficulty, &s.Language, &createdAt); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		s.ID = parsedID
		out = append(out, &s)
	}
	return out, nil
}

func (r *SongRepositoryAdapter) FindByID(id uuid.UUID, lang string) (*domain.Song, error) {
	row := r.db.QueryRow("SELECT id, title_en, title_vi, artist, source, difficulty, language, created_at FROM songs WHERE id = @p1", id)
	var idStr string
	var s domain.Song
	var createdAt interface{}
	if err := row.Scan(&idStr, &s.TitleEN, &s.TitleVI, &s.Artist, &s.Source, &s.Difficulty, &s.Language, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	parsedID, _ := uuid.Parse(idStr)
	s.ID = parsedID
	return &s, nil
}

func (r *SongRepositoryAdapter) Save(song *domain.Song) error {
	_, err := r.db.Exec(`INSERT INTO songs (id, title_en, title_vi, artist, source, difficulty, language, created_at) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`,
		song.ID, song.TitleEN, song.TitleVI, song.Artist, song.Source, song.Difficulty, song.Language, song.CreatedAt)
	return err
}


