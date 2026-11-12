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
	rows, err := r.db.Query(`
		SELECT
			CAST(b.id AS NVARCHAR(50)),
			CAST(b.lesson_id AS NVARCHAR(50)),
			b.block_type,
			b.sort_order,
			b.media_asset_id,
			b.metadata_json,
			b.created_at,
			ISNULL(ten.content, ''),
			ISNULL(tvi.content, '')
		FROM content.lesson_blocks b
		LEFT JOIN content.lesson_block_translations ten ON ten.block_id = b.id AND ten.language = 'EN'
		LEFT JOIN content.lesson_block_translations tvi ON tvi.block_id = b.id AND tvi.language = 'VI'
		WHERE b.lesson_id = @p1
		ORDER BY b.sort_order`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var blocks []*domain.LessonBlock
	for rows.Next() {
		var id, lessonIDStr string
		var blockType string
		var sortOrder int
		var mediaAssetIDRaw interface{}
		var metadataJSON *string
		var createdAt interface{}
		var contentEN, contentVI string
		err := rows.Scan(&id, &lessonIDStr, &blockType, &sortOrder, &mediaAssetIDRaw, &metadataJSON, &createdAt, &contentEN, &contentVI)
		if err != nil {
			return nil, err
		}
		parsedID, _ := uuid.Parse(id)
		parsedLessonID, _ := uuid.Parse(lessonIDStr)
		// media_asset_id is a UUID stored as UNIQUEIDENTIFIER — convert to *string
		var mediaAssetID *string
		if mediaAssetIDRaw != nil {
			if b, ok := mediaAssetIDRaw.([]byte); ok {
				s := uuid.UUID(b).String()
				mediaAssetID = &s
			}
		}
		blocks = append(blocks, &domain.LessonBlock{
			ID:           parsedID,
			LessonID:     parsedLessonID,
			BlockType:    common.ContentBlockType(blockType),
			SortOrder:    sortOrder,
			MediaAssetID: mediaAssetID,
			MetadataJSON: metadataJSON,
			ContentEN:    contentEN,
			ContentVI:    contentVI,
		})
	}
	return blocks, nil
}
