package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/behummble/Questions-answers/internal/models"
)

type Service struct {
	log *slog.Logger
	questionStorage StorageQuestion
	answerStorage StorageAnswer
}

type StorageQuestion interface {
	CreateQuestion(ctx context.Context, data models.Question) error
	Question(ctx context.Context, id int) (models.GetQuestionResponse, error)
	AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error)
	DeleteQuestion(ctx context.Context, id int) error
	Exist(ctx context.Context, id int) (bool, error)
	Shutdown(ctx context.Context)
}

type StorageAnswer interface {
	CreateAnswer(ctx context.Context, data models.Answer) error
	GetAnswer(ctx context.Context, id int) (models.Answer, error)
	DeleteAnswer(ctx context.Context, id int) error
	Shutdown(ctx context.Context)
}

func NewService(log *slog.Logger, questionStorage StorageQuestion, answerStorage StorageAnswer) *Service {
	return &Service{
		log: log,
		questionStorage: questionStorage,
		answerStorage: answerStorage,
	}
}

func(s *Service) Shutdown(ctx context.Context) {
	s.answerStorage.Shutdown(ctx)
	s.questionStorage.Shutdown(ctx)
}

func(s *Service) NewQuestion(ctx context.Context, question []byte) error {
	var questionData models.Question
	err := json.Unmarshal(question, &questionData)
	if err != nil {
		s.log.Error("Json parse error", err)
		return err
	}

	return s.questionStorage.CreateQuestion(ctx, questionData)
}

func(s *Service) Question(ctx context.Context, id int) (models.GetQuestionResponse, error) {	
	return s.questionStorage.Question(ctx, id)
}

func(s *Service) AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error) {
	return s.questionStorage.AllQuestions(ctx)
}

func(s *Service) DeleteQuestion(ctx context.Context, id int) error {
	return s.questionStorage.DeleteQuestion(ctx, id)
}

func(s *Service) NewAnswer(ctx context.Context, answer []byte, questionID int) error {
	exist, err := s.questionStorage.Exist(ctx, questionID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Question with sended id not exist")
	}
	
	var answerData models.Answer
	err = json.Unmarshal(answer, &answerData)
	if err != nil {
		s.log.Error("Json parse error", err)
		return err
	}

	return s.answerStorage.CreateAnswer(ctx, answerData)
}

func(s *Service) Answer(ctx context.Context, id int) (models.Answer, error) {
	return s.answerStorage.GetAnswer(ctx, id)
}

func(s *Service) DeleteAnswer(ctx context.Context, id int) error {
	return s.answerStorage.DeleteAnswer(ctx, id)
}