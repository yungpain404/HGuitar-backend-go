package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/song/domain"
)

type SongResponse struct {
	ID   string `json:"id"`
	Title string `json:"title"`
	Artist string `json:"artist"`
	Source string `json:"source"`
	Difficulty string `json:"difficulty"`
	Language string `json:"language"`
}

type SongSectionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Chords string `json:"chords"`
	Order int `json:"order"`
}

type SongChordResponse struct {
	ID   string `json:"id"`
	ChordID string `json:"chord_id"`
}

type CreateSongRequest struct {
	TitleEN string `json:"title_en" validate:"required"`
	TitleVI string `json:"title_vi"`
	Artist string `json:"artist" validate:"required"`
	Source string `json:"source" validate:"required"`
	Difficulty string `json:"difficulty" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type SongService interface {
	ListSongs(source, difficulty, lang string) ([]*SongResponse, error)
	GetSong(id uuid.UUID, lang string) (*SongResponse, error)
	CreateSong(req CreateSongRequest) (*SongResponse, error)
}

type songService struct {
	repo domain.SongRepository
}

func NewSongService(repo domain.SongRepository) SongService {
	return &songService{repo: repo}
}

func (s *songService) ListSongs(source, difficulty, lang string) ([]*SongResponse, error) {
	songs, err := s.repo.FindAll(source, difficulty, lang)
	if err != nil {
		return nil, err
	}
	resp := make([]*SongResponse, len(songs))
	for i, song := range songs {
		resp[i] = &SongResponse{
			ID:   song.ID.String(),
			Title: song.Title(common.FromCode(lang)),
			Artist: song.Artist,
			Source: song.Source,
			Difficulty: song.Difficulty,
			Language: song.Language,
		}
	}
	return resp, nil
}

func (s *songService) GetSong(id uuid.UUID, lang string) (*SongResponse, error) {
	song, err := s.repo.FindByID(id, lang)
	if err != nil {
		return nil, err
	}
	if song == nil {
		return nil, nil
	}
	return &SongResponse{
		ID:   song.ID.String(),
		Title: song.Title(common.FromCode(lang)),
		Artist: song.Artist,
		Source: song.Source,
		Difficulty: song.Difficulty,
		Language: song.Language,
	}, nil
}

func (s *songService) CreateSong(req CreateSongRequest) (*SongResponse, error) {
	song := &domain.Song{
		ID:   uuid.New(),
		TitleEN: req.TitleEN,
		TitleVI: req.TitleVI,
		Artist: req.Artist,
		Source: req.Source,
		Difficulty: req.Difficulty,
		Language: req.Language,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(song); err != nil {
		return nil, err
	}
	return &SongResponse{
		ID:   song.ID.String(),
		Title: song.TitleEN,
		Artist: song.Artist,
		Source: song.Source,
		Difficulty: song.Difficulty,
		Language: song.Language,
	}, nil
}
