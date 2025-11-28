package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/behummble/Questions-answers/internal/config"
	serv "github.com/behummble/Questions-answers/internal/handlers/http"
	"github.com/behummble/Questions-answers/internal/models"
)

func TestCreateQuestion(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("POST", "/questions", bytes.NewReader([]byte("test")))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
    }
	bodyRes := models.CreateQuestionResponse{}
    res, _ := json.Marshal(bodyRes)
	resStr := string(res)
    if rr.Body.String() != resStr {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), resStr)
    }
}

func TestCreateEmptyQuestion(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("POST", "/questions", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
   s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestGetAllQuestion(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/questions", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	bodyRes := models.GetQuestionsResponse{}
    res, _ := json.Marshal(bodyRes)
	resStr := string(res)
    if rr.Body.String() != resStr {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), resStr)
    }
}

func TestGetQuestion(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/questions/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()

	s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	bodyRes := models.GetQuestionResponse{}
    res, _ := json.Marshal(bodyRes)
	resStr := string(res)
    if rr.Body.String() != resStr {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), resStr)
    }
}

func TestGetQuestionWithoutID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/questions/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}

func TestGetQuestionWithIncorrectID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/questions/abv", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestDeleteQuestion(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("DELETE", "/questions/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNoContent {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNoContent)
    }
}

func TestDeleteQuestionWithoutID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("DELETE", "/questions/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}

func TestDeleteQuestionWithIncorrectID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("DELETE", "/questions/abv", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestCreateCorrectAnswer(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("POST", "/questions/1/answers", bytes.NewReader([]byte("test")))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusCreated)
    }
	bodyRes := models.CreateAnswerResponse{}
    res, _ := json.Marshal(bodyRes)
	resStr := string(res)
    if rr.Body.String() != resStr {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), resStr)
    }
}

func TestCreateAnswerWithoutQuestionID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("POST", "/questions/answers", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusMethodNotAllowed {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}

func TestCreateAnswerWithEmpyBody(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("POST", "/questions/1/answers", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestCreateAnswerWithIncorrectQuestionID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("POST", "/questions/abv/answers", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestGetAnswer(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/answers/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    bodyRes := models.GetAnswerResponse{}
    res, _ := json.Marshal(bodyRes)
	resStr := string(res)
    if rr.Body.String() != resStr {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), resStr)
    }
}

func TestGetAnswerWithoutID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/answers/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}

func TestGetAnswerWithIncorrectID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("GET", "/answers/abv", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func TestDeleteAnswer(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("DELETE", "/answers/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNoContent {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNoContent)
    }
}

func TestDeleteAnswerWithoutID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("DELETE", "/answers/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}

func TestDeleteAnswerWithIncorrectID(t *testing.T) {
	s := createServer()

	req, err := http.NewRequest("DELETE", "/answers/abv", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    s.GetHandler().ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusBadRequest)
    }
}

func createServer() *serv.Server {
	return serv.NewServer(
		context.Background(),
		slog.Default(),
		serverConfig(),
		mockServiceLogic(),
	)
}

func serverConfig() *config.ServerConfig {
	return &config.ServerConfig{}
}

type MockService struct {

}

func mockServiceLogic() *MockService {
	return &MockService{}
}

func(s *MockService) Shutdown(ctx context.Context) {
	
}

func(s *MockService) NewQuestion(ctx context.Context, question []byte) (models.CreateQuestionResponse, error) {
	return models.CreateQuestionResponse{}, nil
}

func(s *MockService) Question(ctx context.Context, id int) (models.GetQuestionResponse, error) {	
	return  models.GetQuestionResponse{}, nil
}

func(s *MockService) AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error) {
	return models.GetQuestionsResponse{}, nil
}

func(s *MockService) DeleteQuestion(ctx context.Context, id int) error {
	return nil
}

func(s *MockService) NewAnswer(ctx context.Context, answer []byte, questionID int) (models.CreateAnswerResponse, error) {
	return models.CreateAnswerResponse{}, nil
}

func(s *MockService) Answer(ctx context.Context, id int) (models.GetAnswerResponse, error) {
	return models.GetAnswerResponse{}, nil
}

func(s *MockService) DeleteAnswer(ctx context.Context, id int) error {
	return nil
}