package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	ldomain "github.com/hguitar/backend-go/internal/learning/domain"
	pdomain "github.com/hguitar/backend-go/internal/progress/domain"
)

type LessonProgressResponse struct {
	LessonID   string    `json:"lesson_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type UserProgressResponse struct {
	CompletedCount int                      `json:"completed_count"`
	Lessons        []*LessonProgressResponse `json:"lessons"`
}

type ProgressService interface {
	CompleteLesson(userID, lessonID uuid.UUID, lang string) (*LessonProgressResponse, error)
	GetUserProgress(userID uuid.UUID, lang string) (*UserProgressResponse, error)
}

type progressService struct {
	progressRepo pdomain.LessonProgressRepository
	lessonRepo   ldomain.LessonRepository
}

func NewProgressService(progressRepo pdomain.LessonProgressRepository, lessonRepo ldomain.LessonRepository) ProgressService {
	return &progressService{progressRepo: progressRepo, lessonRepo: lessonRepo}
}

func (s *progressService) CompleteLesson(userID, lessonID uuid.UUID, lang string) (*LessonProgressResponse, error) {
	_, err := s.lessonRepo.FindByID(lessonID, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	p, err := s.progressRepo.FindByUserIDAndLessonID(userID, lessonID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		p = &pdomain.LessonProgress{
			ID:   uuid.New(),
			UserID: userID,
			LessonID: lessonID,
			Status: "IN_PROGRESS",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	p.Complete()
	if err := s.progressRepo.Save(p); err != nil {
		return nil, err
	}
	return &LessonProgressResponse{
		LessonID:   p.LessonID.String(),
		Status:     p.Status,
		CompletedAt: p.CompletedAt,
	}, nil
}

func (s *progressService) GetUserProgress(userID uuid.UUID, lang string) (*UserProgressResponse, error) {
	progressList, err := s.progressRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	lessonIDs := make([]uuid.UUID, len(progressList))
	for i, p := range progressList {
		lessonIDs[i] = p.LessonID
	}
	lessons, _ := s.lessonRepo.FindAllByIDs(lessonIDs)
	lessonMap := make(map[uuid.UUID]*ldomain.Lesson)
	for _, l := range lessons {
		lessonMap[l.ID] = l
	}
	responses := make([]*LessonProgressResponse, len(progressList))
	completedCount := 0
	for i, p := range progressList {
		title := ""
		if l, ok := lessonMap[p.LessonID]; ok {
			title = l.Title(common.FromCode(lang))
		}
		responses[i] = &LessonProgressResponse{
			LessonID:   p.LessonID.String(),
			Title:      title,
			Status:     p.Status,
			CompletedAt: p.CompletedAt,
		}
		if p.Status == "COMPLETED" {
			completedCount++
		}
	}
	return &UserProgressResponse{
		CompletedCount: completedCount,
		Lessons:        responses,
	}, nil
}
