package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/behummble/Questions-answers/internal/models"
	"github.com/behummble/Questions-answers/internal/mock"
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

func newTestService(questionLen, answerLen int) *Service {
	mockStorageQuestions := mock.NewMockStorageQuestions(questionLen)
	mockStorageAnswers:= mock.NewMockStorageAnswers(answerLen)

	return NewService(
		slog.Default(),
		mockStorageQuestions,
		mockStorageAnswers,
	)
}

func defaultTime() time.Time {
	return time.Date(2000, time.January, 1, 8, 8, 8, 8, time.UTC)
}