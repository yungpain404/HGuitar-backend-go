package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/hguitar/backend-go/internal/assessment/domain"
	"github.com/hguitar/backend-go/internal/common"
)

type AnswerResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type QuestionResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Answers []AnswerResponse `json:"answers"`
}

type QuizResponse struct {
	ID   string `json:"id"`
	Title string `json:"title"`
	Questions []QuestionResponse `json:"questions"`
}

type AttemptResponse struct {
	ID   string `json:"id"`
	QuizID string `json:"quiz_id"`
	Score int `json:"score"`
	Status string `json:"status"`
}

type SubmitAttemptRequest struct {
	Answers map[string]string `json:"answers"`
}

type QuizService interface {
	GetQuizForLesson(lessonID uuid.UUID, lang string) (*QuizResponse, error)
	StartAttempt(quizID, userID uuid.UUID) (*AttemptResponse, error)
	SubmitAttempt(attemptID, userID uuid.UUID, req SubmitAttemptRequest) (*AttemptResponse, error)
	GetAttempt(attemptID, userID uuid.UUID) (*AttemptResponse, error)
}

type quizService struct {
	quizRepo    domain.QuizRepository
	attemptRepo domain.QuizAttemptRepository
}

func NewQuizService(quizRepo domain.QuizRepository, attemptRepo domain.QuizAttemptRepository) QuizService {
	return &quizService{quizRepo: quizRepo, attemptRepo: attemptRepo}
}

func (s *quizService) GetQuizForLesson(lessonID uuid.UUID, lang string) (*QuizResponse, error) {
	quiz, err := s.quizRepo.FindByLessonID(lessonID, common.FromCode(lang))
	if err != nil {
		return nil, err
	}
	if quiz == nil {
		return nil, nil
	}
	return toQuizResponse(quiz), nil
}

func (s *quizService) StartAttempt(quizID, userID uuid.UUID) (*AttemptResponse, error) {
	attempt := &domain.QuizAttempt{
		ID:   uuid.New(),
		QuizID: quizID,
		UserID: userID,
		Score: 0,
		Status: "IN_PROGRESS",
		CreatedAt: time.Now(),
	}
	if err := s.attemptRepo.Create(attempt); err != nil {
		return nil, err
	}
	return &AttemptResponse{
		ID:   attempt.ID.String(),
		QuizID: attempt.QuizID.String(),
		Score: attempt.Score,
		Status: attempt.Status,
	}, nil
}

func (s *quizService) SubmitAttempt(attemptID, userID uuid.UUID, req SubmitAttemptRequest) (*AttemptResponse, error) {
	attempt, err := s.attemptRepo.FindByID(attemptID)
	if err != nil {
		return nil, err
	}
	if attempt == nil {
		return nil, nil
	}
	attempt.Status = "COMPLETED"
	attempt.Score = len(req.Answers)
	if err := s.attemptRepo.Update(attempt); err != nil {
		return nil, err
	}
	return &AttemptResponse{
		ID:   attempt.ID.String(),
		QuizID: attempt.QuizID.String(),
		Score: attempt.Score,
		Status: attempt.Status,
	}, nil
}

func (s *quizService) GetAttempt(attemptID, userID uuid.UUID) (*AttemptResponse, error) {
	attempt, err := s.attemptRepo.FindByID(attemptID)
	if err != nil {
		return nil, err
	}
	if attempt == nil {
		return nil, nil
	}
	return &AttemptResponse{
		ID:   attempt.ID.String(),
		QuizID: attempt.QuizID.String(),
		Score: attempt.Score,
		Status: attempt.Status,
	}, nil
}

func toQuizResponse(q *domain.Quiz) *QuizResponse {
	questions := make([]QuestionResponse, len(q.Questions))
	for i, question := range q.Questions {
		answers := make([]AnswerResponse, len(question.Answers))
		for j, answer := range question.Answers {
			answers[j] = AnswerResponse{ID: answer.ID.String(), Text: answer.Text}
		}
		questions[i] = QuestionResponse{ID: question.ID.String(), Text: question.TextEN, Answers: answers}
	}
	return &QuizResponse{ID: q.ID.String(), Title: q.TitleEN, Questions: questions}
}
