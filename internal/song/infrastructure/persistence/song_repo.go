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
	query := `
		SELECT
			CAST(s.id AS NVARCHAR(50)),
			s.source,
			s.difficulty,
			s.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.artist, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.artist, '')
		FROM song.songs s
		LEFT JOIN song.song_translations ten ON ten.song_id = s.id AND ten.language = 'EN'
		LEFT JOIN song.song_translations tvi ON tvi.song_id = s.id AND tvi.language = 'VI'
		WHERE s.status = 'PUBLISHED'`
	var args []any
	argIdx := 1
	if source != "" {
		query += " AND s.source = @p" + string(rune('0'+argIdx))
		args = append(args, source)
		argIdx++
	}
	if difficulty != "" {
		query += " AND s.difficulty = @p" + string(rune('0'+argIdx))
		args = append(args, difficulty)
	}
	query += " ORDER BY s.sort_order"
	return r.scan(query, args...)
}

func (r *SongRepositoryAdapter) FindByID(id uuid.UUID, lang string) (*domain.Song, error) {
	list, err := r.scan(`
		SELECT
			CAST(s.id AS NVARCHAR(50)),
			s.source,
			s.difficulty,
			s.created_at,
			ISNULL(ten.title, ''),
			ISNULL(ten.artist, ''),
			ISNULL(tvi.title, ''),
			ISNULL(tvi.artist, '')
		FROM song.songs s
		LEFT JOIN song.song_translations ten ON ten.song_id = s.id AND ten.language = 'EN'
		LEFT JOIN song.song_translations tvi ON tvi.song_id = s.id AND tvi.language = 'VI'
		WHERE s.id = @p1`, id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (r *SongRepositoryAdapter) Save(song *domain.Song) error {
	_, err := r.db.Exec(`INSERT INTO song.songs (id, source, difficulty, status, sort_order, created_at, updated_at) VALUES (@p1, @p2, @p3, 'PUBLISHED', 0, GETUTCDATE(), GETUTCDATE())`,
		song.ID, song.Source, song.Difficulty)
	return err
}

func (r *SongRepositoryAdapter) scan(query string, args ...any) ([]*domain.Song, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Song
	for rows.Next() {
		var id, source, difficulty string
		var createdAt interface{}
		var titleEN, artistEN, titleVI, artistVI string
		if err := rows.Scan(&id, &source, &difficulty, &createdAt,
			&titleEN, &artistEN, &titleVI, &artistVI); err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		out = append(out, &domain.Song{
			ID:         parsedID,
			TitleEN:    titleEN,
			TitleVI:    titleVI,
			Artist:     artistEN, // default to EN artist
			Source:     source,
			Difficulty: difficulty,
		})
	}
	return out, nil
}
