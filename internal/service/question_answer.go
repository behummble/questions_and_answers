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
	CreateQuestion(ctx context.Context, data *models.Question) error
	Question(ctx context.Context, id int) (models.QuestionWithAnswers, error)
	AllQuestions(ctx context.Context) ([]models.Question, error)
	DeleteQuestion(ctx context.Context, id int) error
	Exist(ctx context.Context, id int) (bool, error)
	Shutdown(ctx context.Context)
}

type StorageAnswer interface {
	CreateAnswer(ctx context.Context, data *models.Answer) error
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

func(s *Service) NewQuestion(ctx context.Context, question []byte) (models.CreateQuestionResponse, error) {
	var questionRequest models.CreateQuestionRequest
	err := json.Unmarshal(question, &questionRequest)
	if err != nil {
		s.log.Error(
			"ParsingJSONError", 
			slog.String("component", "json/unmarshalling"),
			slog.Any("error", err),
		)
		return models.CreateQuestionResponse{}, errors.New("DecodingDataError")
	}

	if questionRequest.Text == "" {
		return models.CreateQuestionResponse{}, errors.New("BodyExecutionError")
	}

	questionData := models.Question{
		Text: questionRequest.Text,
	}

	err = s.questionStorage.CreateQuestion(ctx, &questionData)
	if err != nil {
		s.log.Error(
			"DB_WritingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return models.CreateQuestionResponse{}, errors.New("DB_WritingError")
	}

	return models.CreateQuestionResponse{Question: questionData}, err
}

func(s *Service) Question(ctx context.Context, id int) (models.GetQuestionResponse, error) {	
	res, err := s.questionStorage.Question(ctx, id)
	if err != nil {
		s.log.Error(
			"DB_ReadingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return  models.GetQuestionResponse{}, err
	}

	return models.GetQuestionResponse{Question: res.Question, Answers: res.Answers}, nil
}

func(s *Service) AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error) {
	allQuestions, err := s.questionStorage.AllQuestions(ctx)
	if err != nil {
		s.log.Error(
			"DB_ReadingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return  models.GetQuestionsResponse{}, err
	}

	return models.GetQuestionsResponse{Questions: allQuestions}, err
}

func(s *Service) DeleteQuestion(ctx context.Context, id int) error {
	err := s.questionStorage.DeleteQuestion(ctx, id)
	if err != nil {
		s.log.Error(
			"DB_DeletingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return errors.New("DB_DeletingError")
	}
	return nil
}

func(s *Service) NewAnswer(ctx context.Context, answer []byte, questionID int) (models.CreateAnswerResponse, error) {
	exist, err := s.questionStorage.Exist(ctx, questionID)
	if err != nil {
		s.log.Error(
			"DB_ReadingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return models.CreateAnswerResponse{}, errors.New("DBReadingError")
	}
	if !exist {
		return models.CreateAnswerResponse{}, errors.New("QuestionWithSendedIDNotExist")
	}
	
	var answerRequest models.CreateAnswerRequest
	err = json.Unmarshal(answer, &answerRequest)
	if err != nil {
		s.log.Error(
			"ParsingJSONError", 
			slog.String("component", "json/unmarshaling"),
			slog.Any("error", err),
		)
		return models.CreateAnswerResponse{}, errors.New("ParsingJSONError")
	}

	if answerRequest.Text == "" {
		return models.CreateAnswerResponse{}, errors.New("BodyExecutionError")
	}

	answerData := models.Answer{
		Text: answerRequest.Text,
		UserID: answerRequest.UserID,
		QuestionID: questionID,
	}

	err = s.answerStorage.CreateAnswer(ctx, &answerData)
	if err != nil {
		s.log.Error(
			"DB_WritingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return models.CreateAnswerResponse{}, errors.New("DB_WritingError")
	}

	return models.CreateAnswerResponse{Answer: answerData}, err
}

func(s *Service) Answer(ctx context.Context, id int) (models.GetAnswerResponse, error) {
	answer, err := s.answerStorage.GetAnswer(ctx, id)
	if err != nil {
		s.log.Error(
			"DB_ReadingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
	}
	return models.GetAnswerResponse{Answer: answer}, err
}

func(s *Service) DeleteAnswer(ctx context.Context, id int) error {
	err := s.answerStorage.DeleteAnswer(ctx, id)
	if err != nil {
		s.log.Error(
			"DB_DeletingError", 
			slog.String("component", "db"),
			slog.Any("error", err),
		)
		return errors.New("DB_DeletingError")
	}

	return nil
}