package application

import (
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/content/domain"
)

type LessonBlockResponse struct {
	ID         string                `json:"id"`
	LessonID   string                `json:"lesson_id"`
	BlockType  string                `json:"block_type"`
	SortOrder  int                   `json:"sort_order"`
	MediaAssetID *string              `json:"media_asset_id,omitempty"`
	MetadataJSON *string              `json:"metadata_json,omitempty"`
	Content    string                `json:"content"`
}

type ContentService interface {
	GetBlocksForLesson(lessonID uuid.UUID, lang string) ([]*LessonBlockResponse, error)
}

type contentService struct {
	repo domain.LessonBlockRepository
}

func NewContentService(repo domain.LessonBlockRepository) ContentService {
	return &contentService{repo: repo}
}

func (s *contentService) GetBlocksForLesson(lessonID uuid.UUID, lang string) ([]*LessonBlockResponse, error) {
	blocks, err := s.repo.FindByLessonIDOrderBySortOrder(lessonID)
	if err != nil {
		return nil, err
	}
	language := common.FromCode(lang)
	resp := make([]*LessonBlockResponse, len(blocks))
	for i, b := range blocks {
		resp[i] = &LessonBlockResponse{
			ID:          b.ID.String(),
			LessonID:    b.LessonID.String(),
			BlockType:   string(b.BlockType),
			SortOrder:   b.SortOrder,
			MediaAssetID: b.MediaAssetID,
			MetadataJSON: b.MetadataJSON,
			Content:     b.Content(language),
		}
	}
	return resp, nil
}
