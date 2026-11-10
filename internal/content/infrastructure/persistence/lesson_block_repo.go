package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/content/domain"
	"github.com/hguitar/backend-go/internal/common"
)

type LessonBlockRepositoryAdapter struct {
	db *sql.DB
}

func NewLessonBlockRepositoryAdapter(db *sql.DB) *LessonBlockRepositoryAdapter {
	return &LessonBlockRepositoryAdapter{db: db}
}

func (r *LessonBlockRepositoryAdapter) FindByLessonIDOrderBySortOrder(lessonID uuid.UUID) ([]*domain.LessonBlock, error) {
	rows, err := r.db.Query("SELECT id, lesson_id, block_type, sort_order, media_asset_id, metadata_json, content_en, content_vi, created_at, updated_at FROM lesson_blocks WHERE lesson_id = @p1 ORDER BY sort_order", lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var blocks []*domain.LessonBlock
	for rows.Next() {
		var id, lessonID string
		var blockType string
		var sortOrder int
		var mediaAssetID, metadataJSON *string
		var contentEN, contentVI string
		var createdAt, updatedAt interface{}
		err := rows.Scan(&id, &lessonID, &blockType, &sortOrder, &mediaAssetID, &metadataJSON, &contentEN, &contentVI, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedLessonID, _ := uuid.Parse(lessonID)
		blocks = append(blocks, &domain.LessonBlock{
			ID:          parsedID,
			LessonID:    parsedLessonID,
			BlockType:   common.ContentBlockType(blockType),
			SortOrder:   sortOrder,
			MediaAssetID: mediaAssetID,
			MetadataJSON: metadataJSON,
			ContentEN:   contentEN,
			ContentVI:   contentVI,
		})
	}
	return blocks, nil
}
