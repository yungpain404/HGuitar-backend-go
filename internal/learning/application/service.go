package application

import (
	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/common"
	"github.com/hguitar/backend-go/internal/learning/domain"
)

type CourseResponse struct {
	ID   string `json:"id"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	GuitarType string `json:"guitar_type"`
	Level string `json:"level"`
}

type ModuleResponse struct {
	ID   string `json:"id"`
	CourseID string `json:"course_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	SortOrder int `json:"sort_order"`
}

type LessonResponse struct {
	ID   string `json:"id"`
	ModuleID string `json:"module_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	VideoURL string `json:"video_url"`
	SortOrder int `json:"sort_order"`
	DurationMin int `json:"duration_min"`
}

type LearningMethodResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Difficulty string `json:"difficulty"`
}

type LearningQueryService interface {
	GetCourses(lang string, guitarType string) ([]*CourseResponse, error)
	GetCourse(id uuid.UUID, lang string) (*CourseResponse, error)
	GetModulesByCourse(id uuid.UUID, lang string) ([]*ModuleResponse, error)
	GetModule(id uuid.UUID, lang string) (*ModuleResponse, error)
	GetLessonsByModule(moduleID uuid.UUID, lang string) ([]*LessonResponse, error)
	GetLesson(id uuid.UUID, lang string) (*LessonResponse, error)
}

type LearningMethodService interface {
	ListMethods(difficulty string, lang string) ([]*LearningMethodResponse, error)
	GetMethod(id uuid.UUID, lang string) (*LearningMethodResponse, error)
}

type learningQueryService struct {
	courseRepo  domain.CourseRepository
	moduleRepo  domain.ModuleRepository
	lessonRepo  domain.LessonRepository
}

func NewLearningQueryService(courseRepo domain.CourseRepository, moduleRepo domain.ModuleRepository, lessonRepo domain.LessonRepository) LearningQueryService {
	return &learningQueryService{courseRepo: courseRepo, moduleRepo: moduleRepo, lessonRepo: lessonRepo}
}

func (s *learningQueryService) GetCourses(lang string, guitarType string) ([]*CourseResponse, error) {
	courses, err := s.courseRepo.FindAll(common.FromCode(lang), guitarType)
	if err != nil {
		return nil, err
	}
	resp := make([]*CourseResponse, len(courses))
	for i, c := range courses {
		resp[i] = &CourseResponse{
			ID:   c.ID.String(),
			Title: c.Title(common.FromCode(lang)),
			Slug: c.Slug,
			Description: c.Description(common.FromCode(lang)),
			GuitarType: c.GuitarType,
			Level: c.Level,
		}
	}
	return resp, nil
}

func (s *learningQueryService) GetCourse(id uuid.UUID, lang string) (*CourseResponse, error) {
	c, err := s.courseRepo.FindByID(id, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return &CourseResponse{
		ID:   c.ID.String(),
		Title: c.Title(common.FromCode(lang)),
		Slug: c.Slug,
		Description: c.Description(common.FromCode(lang)),
		GuitarType: c.GuitarType,
		Level: c.Level,
	}, nil
}

func (s *learningQueryService) GetModulesByCourse(id uuid.UUID, lang string) ([]*ModuleResponse, error) {
	modules, err := s.moduleRepo.FindByCourseID(id, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	resp := make([]*ModuleResponse, len(modules))
	for i, m := range modules {
		resp[i] = &ModuleResponse{
			ID:   m.ID.String(),
			CourseID: m.CourseID.String(),
			Title: m.Title(common.FromCode(lang)),
			Description: m.DescriptionEN,
			SortOrder: m.SortOrder,
		}
	}
	return resp, nil
}

func (s *learningQueryService) GetModule(id uuid.UUID, lang string) (*ModuleResponse, error) {
	m, err := s.moduleRepo.FindByID(id, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return &ModuleResponse{
		ID:   m.ID.String(),
		CourseID: m.CourseID.String(),
		Title: m.Title(common.FromCode(lang)),
		Description: m.DescriptionEN,
		SortOrder: m.SortOrder,
	}, nil
}

func (s *learningQueryService) GetLessonsByModule(moduleID uuid.UUID, lang string) ([]*LessonResponse, error) {
	lessons, err := s.lessonRepo.FindByModuleID(moduleID, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	resp := make([]*LessonResponse, len(lessons))
	for i, l := range lessons {
		resp[i] = toLessonResponse(l, common.FromCode(lang))
	}
	return resp, nil
}

func (s *learningQueryService) GetLesson(id uuid.UUID, lang string) (*LessonResponse, error) {
	l, err := s.lessonRepo.FindByID(id, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, nil
	}
	return toLessonResponse(l, common.FromCode(lang)), nil
}

func toLessonResponse(l *domain.Lesson, lang common.Language) *LessonResponse {
	return &LessonResponse{
		ID:   l.ID.String(),
		ModuleID: l.ModuleID.String(),
		Title: l.Title(lang),
		Content: l.ContentEN, // đơn giản
		VideoURL: l.VideoURL,
		SortOrder: l.SortOrder,
		DurationMin: l.DurationMin,
	}
}

type learningMethodService struct {
	repo domain.LearningMethodRepository
}

func NewLearningMethodService(repo domain.LearningMethodRepository) LearningMethodService {
	return &learningMethodService{repo: repo}
}

func (s *learningMethodService) ListMethods(difficulty string, lang string) ([]*LearningMethodResponse, error) {
	methods, err := s.repo.FindAll(difficulty, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	resp := make([]*LearningMethodResponse, len(methods))
	for i, m := range methods {
		resp[i] = &LearningMethodResponse{
			ID:   m.ID.String(),
			Name: m.Name(common.FromCode(lang)),
			Description: m.DescriptionEN,
			Difficulty: m.Difficulty,
		}
	}
	return resp, nil
}

func (s *learningMethodService) GetMethod(id uuid.UUID, lang string) (*LearningMethodResponse, error) {
	m, err := s.repo.FindByID(id, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return &LearningMethodResponse{
		ID:   m.ID.String(),
		Name: m.Name(common.FromCode(lang)),
		Description: m.DescriptionEN,
		Difficulty: m.Difficulty,
	}, nil
}
