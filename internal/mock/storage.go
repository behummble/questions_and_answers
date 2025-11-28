package mock

import (
	"time"
	"context"
	"errors"
	"github.com/behummble/Questions-answers/internal/models"
)

type MockStorageQuestions struct {
	db map[int]models.Question
}

type MockStorageAnswers struct {
	db map[int]models.Answer
}

func NewMockStorageAnswers(len int) *MockStorageAnswers {
	return &MockStorageAnswers{
		db: make(map[int]models.Answer, len),
	}
}

func NewMockStorageQuestions(len int) *MockStorageQuestions {
	return &MockStorageQuestions{
		db: make(map[int]models.Question, len),
	}
}

func(s *MockStorageAnswers) CreateAnswer(ctx context.Context, data *models.Answer) error {
	ind := len(s.db) + 1
	data.CreatedAt = defaultTime()
	data.ID = ind
	s.db[ind] = *data
	return nil
}

func(s *MockStorageAnswers) GetAnswer(ctx context.Context, id int) (models.Answer, error) {
	res, ok := s.db[id]
	if !ok {
		return res, errors.New("NotFound")
	}

	return res, nil
}

func(s *MockStorageAnswers) DeleteAnswer(ctx context.Context, id int) (int, error){
	if _, ok := s.db[id]; !ok {
		return 0, errors.New("NotFound")
	}
	delete(s.db, id)
	return 1, nil
}

func(s *MockStorageQuestions) CreateQuestion(ctx context.Context, data *models.Question) error {
	ind := len(s.db) + 1
	data.CreatedAt = defaultTime()
	data.ID = ind
	s.db[ind] = *data
	return nil
}

func(s *MockStorageQuestions) Question(ctx context.Context, id int) (models.QuestionWithAnswers, error) {
	res, ok := s.db[id]
	if !ok {
		return models.QuestionWithAnswers{}, errors.New("NotFound")
	}

	return models.QuestionWithAnswers{Question: res}, nil
}

func(s *MockStorageQuestions) AllQuestions(ctx context.Context) ([]models.Question, error) {
	res := make([]models.Question, 0, len(s.db))
	for _, v := range s.db {
		res = append(res, v)
	}

	return res, nil
}

func(s *MockStorageQuestions) DeleteQuestion(ctx context.Context, id int) (int, error) {
	if _, ok := s.db[id]; !ok {
		return 0, errors.New("NotFound")
	}
	delete(s.db, id)
	return 1, nil
}

func(s *MockStorageQuestions) Exist(ctx context.Context, id int) (models.Question, error) {
	_, ok := s.db[id]
	
	if !ok {
		return models.Question{}, errors.New("NotFound")
	}
	return models.Question{}, nil
}

func(s *MockStorageAnswers) AllAnswers(questionID int) []models.Answer {
	res := make([]models.Answer, 0, len(s.db))
	for l, v := range s.db {
		if l == questionID {
			res = append(res, v)
		}
	}

	return res
}

func(s *MockStorageAnswers) DeleteAllAnswers(questionID int) {
	del := make([]int, len(s.db))
	for l := range s.db {
		if l == questionID {
			del = append(del, l)
		}
	}
	for _, v := range del {
		delete(s.db, v)
	}
}

func(s *MockStorageAnswers) Shutdown(ctx context.Context) {

}

func(s *MockStorageQuestions) Shutdown(ctx context.Context) {
	
}

func defaultTime() time.Time {
	return time.Date(2000, time.January, 1, 8, 8, 8, 8, time.UTC)
}