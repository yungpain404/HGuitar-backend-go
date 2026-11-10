package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
)

type LessonBlock struct {
	ID          uuid.UUID
	LessonID    uuid.UUID
	BlockType   common.ContentBlockType
	SortOrder   int
	MediaAssetID *string
	MetadataJSON *string
	ContentEN   string
	ContentVI   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *LessonBlock) Content(lang common.Language) string {
	if lang == common.LanguageVI {
		return b.ContentVI
	}
	return b.ContentEN
}

type LessonBlockRepository interface {
	FindByLessonIDOrderBySortOrder(lessonID uuid.UUID) ([]*LessonBlock, error)
}
