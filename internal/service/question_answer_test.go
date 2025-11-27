package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/behummble/Questions-answers/internal/models"
)

func TestNewQuestionCorrect(t *testing.T) {
	service := newTestService(1, 1)
	excpected := models.CreateQuestionResponse {
		Question: models.Question{
			ID: 1,
			Text: "test",
			CreatedAt: defaultTime(),
		},
	}

	res, err := CreateQuestion(service, t)
	if err != nil {
        t.Fatal(err)
    }

	if res != excpected {
		t.Errorf("Excpected result %+v, got %+v", excpected, res)
	}
}

func TestNewQuestionWithInvalidData(t *testing.T) {
	service := newTestService(1, 1)

	raw := "{\"struct\":\"\"}"

	_, err := service.NewQuestion(context.Background(), []byte(raw))
	if err == nil {
        t.Error("Excpected error")
    }
}

func TestQuestionWithCorrectID(t *testing.T) {
	service := newTestService(1, 1)

	created, err := CreateQuestion(service, t)
	if err != nil {
        t.Fatal(err)
    }

	_, err = service.Question(context.Background(), created.Question.ID)
	if err != nil {
        t.Error("Excpect value, not error")
    }
}

func TestQuestionWithInvalidID(t *testing.T) {
	service := newTestService(1, 1)

	_, err := service.Question(context.Background(), 2)
	if err == nil {
        t.Fatal("Excpect error")
    }
}

func TestAllQuestionsExists(t *testing.T) {
	service := newTestService(3, 1)

	for i := 0; i < 3; i++ {
		_, err := CreateQuestion(service, t)
		if err != nil {
			t.Fatal(err)
		}
	}

	all, err := service.AllQuestions(context.Background())
	if err != nil {
        t.Fatal("Unexcpected error")
    }

	if len(all.Questions) != 3 {
		t.Error("Excpect len 3 of questions")
	}
}

func TestAllQuestionsEmpty(t *testing.T) {
	service := newTestService(3, 1)

	all, err := service.AllQuestions(context.Background())
	if err != nil {
        t.Fatal("Unexcpected error")
    }

	if len(all.Questions) != 0 {
		t.Error("Excpect empty array")
	}
}

func TestDeleteQuestionCorrect(t *testing.T) {
	service := newTestService(1, 1)

	created, err := CreateQuestion(service, t)
	if err != nil {
        t.Fatal("Unexcpected error")
    }

	err = service.DeleteQuestion(context.Background(), created.Question.ID)
	if err != nil {
        t.Fatal("Unexcpected error")
    }
}

func TestDeleteQuestionWithInvalidID(t *testing.T) {
	service := newTestService(1, 1)

	err := service.DeleteQuestion(context.Background(), 1)
	if err == nil {
        t.Error("Unexcpected error")
    }
}

func TestNewAnswerCorrect(t *testing.T) {
	service := newTestService(1, 1)
	question, err := CreateQuestion(service, t)
	if err != nil {
        t.Fatal("Unexcpected error")
    }

	_, err = service.questionStorage.Question(context.Background(), question.Question.ID)
	if err != nil {
		t.Fatal(err)
	}

	res, err := CreateAnswer(service, question.Question.ID, t)
	if err != nil {
        t.Error("Unexcpected Error")
    }

	excpected := models.CreateAnswerResponse {
		Answer: models.Answer{
			ID: 1,
			Text: "test",
			CreatedAt: time.Date(2000, time.January, 1, 8, 8, 8, 8, time.UTC),
			UserID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			QuestionID: question.Question.ID,
		},
	}

	if res != excpected {
		t.Errorf("Excpected result %+v, got %+v", excpected, res)
	}

}

func TestNewAnswerWithInvalidQuestionID(t *testing.T) {
	service := newTestService(1, 1)

	_, err := service.questionStorage.Question(context.Background(), 2)
	if err == nil {
		t.Error("Excepcted error")
	}
}

func TestAnswer(t *testing.T) {
	service := newTestService(1, 1)

	_, err := CreateQuestion(service, t)
	if err != nil {
        t.Fatal("Unexcpected Error")
    }

	created, err := CreateAnswer(service, 1, t)
	if err != nil {
        t.Fatal("Unexcpected Error")
    }

	res, err := service.Answer(context.Background(), 1)
	if err != nil {
        t.Error("Unexcpected Error")
    }

	excpected := models.GetAnswerResponse{Answer: created.Answer}

	if res != excpected {
		t.Errorf("Excpected result %+v, got %+v", created, res)
	}
}

func TestAnswerWithInvalidID(t *testing.T) {
	service := newTestService(1, 1)
	_, err := service.Answer(context.Background(), 1)

	if err == nil {
        t.Error("Excpected Error")
    }
}

func TestDeleteAnswer(t *testing.T) {
	service := newTestService(1, 1)
	_, err := CreateQuestion(service, t)
	if err != nil {
        t.Fatal("Unexcpected Error")
    }
	_, err = CreateAnswer(service, 1, t)
	if err != nil {
        t.Fatal("Unexcpected Error")
    }

	err = service.DeleteAnswer(context.Background(), 1)

	if err != nil {
        t.Error("Unexcpected Error")
    }
}

func TestDeleteAnswerWithInvalidID(t *testing.T) {
	service := newTestService(1, 1)

	err := service.DeleteAnswer(context.Background(), 1)

	if err == nil {
        t.Error("Excpected Error")
    }
}

func CreateQuestion(service *Service, t *testing.T) (models.CreateQuestionResponse, error) {
	m := models.CreateQuestionRequest{
		Text: "test",
	}
	raw, err := json.Marshal(m)
	if err != nil {
        t.Fatal(err)
    }
	res, err := service.NewQuestion(context.Background(), raw)
	if err != nil {
        t.Fatal(err)
    }
	return res, nil
}

func CreateAnswer(service *Service, questionID int, t *testing.T) (models.CreateAnswerResponse, error) {
	m := models.CreateAnswerRequest{
		Text: "test",
		UserID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
	}
	raw, err := json.Marshal(m)
	if err != nil {
        t.Fatal(err)
    }
	res, err := service.NewAnswer(context.Background(), raw, questionID)
	if err != nil {
        t.Fatal(err)
    }
	return res, nil
}

type StorageMockQuestions struct {
	db map[int]models.Question
}

type StorageMockAnswers struct {
	db map[int]models.Answer
}

func defaultTime() time.Time {
	return time.Date(2000, time.January, 1, 8, 8, 8, 8, time.UTC)
}

func(s *StorageMockAnswers) CreateAnswer(ctx context.Context, data *models.Answer) error {
	ind := len(s.db) + 1
	data.CreatedAt = defaultTime()
	data.ID = ind
	s.db[ind] = *data
	return nil
}

func(s *StorageMockAnswers) GetAnswer(ctx context.Context, id int) (models.Answer, error) {
	res, ok := s.db[id]
	if !ok {
		return res, errors.New("NotFound")
	}

	return res, nil
}

func(s *StorageMockAnswers) DeleteAnswer(ctx context.Context, id int) (int, error){
	if _, ok := s.db[id]; !ok {
		return 0, errors.New("NotFound")
	}
	delete(s.db, id)
	return 1, nil
}

func(s *StorageMockQuestions) CreateQuestion(ctx context.Context, data *models.Question) error {
	ind := len(s.db) + 1
	data.CreatedAt = defaultTime()
	data.ID = ind
	s.db[ind] = *data
	return nil
}

func(s *StorageMockQuestions) Question(ctx context.Context, id int) (models.QuestionWithAnswers, error) {
	res, ok := s.db[id]
	if !ok {
		return models.QuestionWithAnswers{}, errors.New("NotFound")
	}

	return models.QuestionWithAnswers{Question: res}, nil
}

func(s *StorageMockQuestions) AllQuestions(ctx context.Context) ([]models.Question, error) {
	res := make([]models.Question, 0, len(s.db))
	for _, v := range s.db {
		res = append(res, v)
	}

	return res, nil
}

func(s *StorageMockQuestions) DeleteQuestion(ctx context.Context, id int) (int, error) {
	if _, ok := s.db[id]; !ok {
		return 0, errors.New("NotFound")
	}
	delete(s.db, id)
	return 1, nil
}

func(s *StorageMockQuestions) Exist(ctx context.Context, id int) (bool, error) {
	_, ok := s.db[id]
	return ok, nil
}

func(s *StorageMockAnswers) AllAnswers(questionID int) []models.Answer {
	res := make([]models.Answer, 0, len(s.db))
	for l, v := range s.db {
		if l == questionID {
			res = append(res, v)
		}
	}

	return res
}

func(s *StorageMockAnswers) DeleteAllAnswers(questionID int) {
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

func(s *StorageMockAnswers) Shutdown(ctx context.Context) {

}

func(s *StorageMockQuestions) Shutdown(ctx context.Context) {
	
}

func newMockStorageQuestions(len int) *StorageMockQuestions {
	return &StorageMockQuestions{
		db: make(map[int]models.Question, len),
	}
}

func newMockStorageAnswers(len int) *StorageMockAnswers {
	return &StorageMockAnswers{
		db: make(map[int]models.Answer, len),
	}
}

func newTestService(questionLen, answerLen int) *Service {
	mockStorageQuestions := newMockStorageQuestions(questionLen)
	mockStorageAnswers:= newMockStorageAnswers(answerLen)

	return NewService(
		slog.Default(),
		mockStorageQuestions,
		mockStorageAnswers,
	)
}