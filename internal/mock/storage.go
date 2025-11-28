package mock

import (
	"time"
	"context"
	"github.com/behummble/Questions-answers/internal/models"
	"gorm.io/gorm"
)

type MockStorageQuestions struct {
	db map[int]models.Question
	id int
	storageAnswers *MockStorageAnswers
}

type MockStorageAnswers struct {
	db map[int]models.Answer
	id int
}

func NewMockStorageAnswers(len int) *MockStorageAnswers {
	return &MockStorageAnswers{
		db: make(map[int]models.Answer, len),
	}
}

func NewMockStorageQuestions(len int, storageAnswers *MockStorageAnswers) *MockStorageQuestions {
	return &MockStorageQuestions{
		db: make(map[int]models.Question, len),
		storageAnswers: storageAnswers,
	}
}

func(s *MockStorageAnswers) CreateAnswer(ctx context.Context, data []*models.Answer) error {
	for _, v := range data {
		s.id += 1
		ind := s.id
		v.CreatedAt = defaultTime()
		v.ID = ind
		s.db[ind] = *v
	}
	
	return nil
}

func(s *MockStorageAnswers) GetAnswer(ctx context.Context, id int) (models.Answer, error) {
	res, ok := s.db[id]
	if !ok {
		return res, gorm.ErrRecordNotFound
	}

	return res, nil
}

func(s *MockStorageAnswers) DeleteAnswer(ctx context.Context, id int) (int, error){
	if _, ok := s.db[id]; !ok {
		return 0, gorm.ErrRecordNotFound
	}
	delete(s.db, id)
	return 1, nil
}

func(s *MockStorageQuestions) CreateQuestion(ctx context.Context, data *models.Question) error {
	s.id += 1
	ind := s.id
	data.CreatedAt = defaultTime()
	data.ID = ind
	s.db[ind] = *data
	return nil
}

func(s *MockStorageQuestions) Question(ctx context.Context, id int) (models.QuestionWithAnswers, error) {
	res, ok := s.db[id]
	if !ok {
		return models.QuestionWithAnswers{}, gorm.ErrRecordNotFound
	}

	answers := s.storageAnswers.AllAnswers(id)

	return models.QuestionWithAnswers{Question: res, Answers: answers}, nil
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
		return 0, gorm.ErrRecordNotFound
	}
	delete(s.db, id)
	s.storageAnswers.DeleteAllAnswers(id)
	return 1, nil
}

func(s *MockStorageQuestions) Exist(ctx context.Context, id int) (models.Question, error) {
	_, ok := s.db[id]
	
	if !ok {
		return models.Question{}, gorm.ErrRecordNotFound
	}
	return models.Question{}, nil
}

func(s *MockStorageAnswers) AllAnswers(questionID int) []models.Answer {
	res := make([]models.Answer, 0, len(s.db))
	for _, v := range s.db {
		if v.QuestionID == questionID {
			res = append(res, v)
		}
	}

	return res
}

func(s *MockStorageAnswers) DeleteAllAnswers(questionID int) {
	del := make([]int, len(s.db))
	for ind, v := range s.db {
		if v.QuestionID == questionID {
			del = append(del, ind)
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