package endtoend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/behummble/Questions-answers/internal/mock"
	"github.com/behummble/Questions-answers/internal/models"
	srv "github.com/behummble/Questions-answers/internal/handlers/http"
	"github.com/behummble/Questions-answers/internal/service"
	"github.com/behummble/Questions-answers/internal/config"
)

func TestEndToEnd(t *testing.T) {
    // Инициализация моковых хранилищ
	mockAnswerStorage := mock.NewMockStorageAnswers(10)
    mockQuestionStorage := mock.NewMockStorageQuestions(10, mockAnswerStorage)
	ctx := context.Background()

    // Инициализация сервиса
    svc := service.NewService(slog.Default(), mockQuestionStorage, mockAnswerStorage)

    // Инициализация сервера
    srv := srv.NewServer(ctx, slog.Default(), &config.ServerConfig{}, svc)
    testServer := httptest.NewServer(srv.GetHandler())
    defer testServer.Close()

    client := testServer.Client()
    
    t.Run("Complete question-answer flow", func(t *testing.T) {
		
        createQuestionReq := models.CreateQuestionRequest{
            Text: "What is Golang?",
        }
        
        reqBody, _ := json.Marshal(createQuestionReq)
        resp, err := client.Post(testServer.URL+"/questions", "application/json", bytes.NewBuffer(reqBody))
        if err != nil {
            t.Fatalf("Failed to create question: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusCreated {
            t.Errorf("Expected status 201, got %d", resp.StatusCode)
        }

        var createQuestionResp models.CreateQuestionResponse
        if err := json.NewDecoder(resp.Body).Decode(&createQuestionResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        questionID := createQuestionResp.Question.ID
        if questionID == 0 {
            t.Error("Expected question ID to be set")
        }

        resp, err = client.Get(testServer.URL + "/questions")
        if err != nil {
            t.Fatalf("Failed to get questions: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status 200, got %d", resp.StatusCode)
        }

        var getQuestionsResp models.GetQuestionsResponse
        if err := json.NewDecoder(resp.Body).Decode(&getQuestionsResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if len(getQuestionsResp.Questions) != 1 {
            t.Errorf("Expected 1 question, got %d", len(getQuestionsResp.Questions))
        }

        // 3. Получение конкретного вопроса
        resp, err = client.Get(fmt.Sprintf("%s/questions/%d", testServer.URL, questionID))
        if err != nil {
            t.Fatalf("Failed to get question: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status 200, got %d", resp.StatusCode)
        }

        var getQuestionResp models.GetQuestionResponse
        if err := json.NewDecoder(resp.Body).Decode(&getQuestionResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if getQuestionResp.Question.ID != questionID {
            t.Errorf("Expected question ID %d, got %d", questionID, getQuestionResp.Question.ID)
        }

        // 4. Создание ответа
        createAnswerReq := models.CreateAnswerRequest{
            Texts:   []string{"Golang is a programming language created by Google"},
            UserID: "123e4567-e89b-12d3-a456-426614174000",
        }
        
        reqBody, _ = json.Marshal(createAnswerReq)
        resp, err = client.Post(
            fmt.Sprintf("%s/questions/%d/answers", testServer.URL, questionID),
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        if err != nil {
            t.Fatalf("Failed to create answer: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusCreated {
            t.Errorf("Expected status 201, got %d", resp.StatusCode)
        }

        var createAnswerResp models.CreateAnswerResponse
        if err := json.NewDecoder(resp.Body).Decode(&createAnswerResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        answerID := createAnswerResp.Answers[0].ID
        if answerID == 0 {
            t.Error("Expected answer ID to be set")
        }

        // 5. Получение ответа
        resp, err = client.Get(fmt.Sprintf("%s/answers/%d", testServer.URL, answerID))
        if err != nil {
            t.Fatalf("Failed to get answer: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("Expected status 200, got %d", resp.StatusCode)
        }

        var getAnswerResp models.GetAnswerResponse
        if err := json.NewDecoder(resp.Body).Decode(&getAnswerResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if getAnswerResp.Answer.ID != answerID {
            t.Errorf("Expected answer ID %d, got %d", answerID, getAnswerResp.Answer.ID)
        }

        // 6. Проверка, что ответ привязан к вопросу
        resp, err = client.Get(fmt.Sprintf("%s/questions/%d", testServer.URL, questionID))
        if err != nil {
            t.Fatalf("Failed to get question with answers: %v", err)
        }
        defer resp.Body.Close()

        if err := json.NewDecoder(resp.Body).Decode(&getQuestionResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if len(getQuestionResp.Answers) != 1 {
            t.Errorf("Expected 1 answer, got %d", len(getQuestionResp.Answers))
        }

		// 7. Удаление ответа
        req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/answers/%d", testServer.URL, answerID), nil)
        if err != nil {
            t.Fatalf("Failed to create delete request: %v", err)
        }

        resp, err = client.Do(req)
        if err != nil {
            t.Fatalf("Failed to delete answer: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNoContent {
            t.Errorf("Expected status 204, got %d", resp.StatusCode)
        }

		// 8. Создание ответа с пустым массивом ответов
		createAnswerReq = models.CreateAnswerRequest{
            Texts:   []string{},
            UserID: "123e4567-e89b-12d3-a456-426614174000",
        }
        
        reqBody, _ = json.Marshal(createAnswerReq)
        resp, err = client.Post(
            fmt.Sprintf("%s/questions/%d/answers", testServer.URL, questionID),
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        if err != nil {
            t.Fatalf("Failed to create answer: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusInternalServerError {
            t.Errorf("Expected status 500, got %d", resp.StatusCode)
        }

        // 9. Создание ответов по несуществующему вопросу
		createAnswerReq = models.CreateAnswerRequest{
            Texts:   []string{"answer"},
            UserID: "123e4567-e89b-12d3-a456-426614174000",
        }
        
        reqBody, _ = json.Marshal(createAnswerReq)
        resp, err = client.Post(
            fmt.Sprintf("%s/questions/%d/answers", testServer.URL, 55),
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        if err != nil {
            t.Fatalf("Failed to create answer: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNotFound {
            t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}

		// 10. Создание нескольких ответов за раз
        createAnswerReq = models.CreateAnswerRequest{
            Texts:   []string{"answer1", "answer2", "answer3"},
            UserID: "123e4567-e89b-12d3-a456-426614174000",
        }
        
        reqBody, _ = json.Marshal(createAnswerReq)
        resp, err = client.Post(
            fmt.Sprintf("%s/questions/%d/answers", testServer.URL, questionID),
            "application/json",
            bytes.NewBuffer(reqBody),
        )
        if err != nil {
            t.Fatalf("Failed to create answer: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusCreated {
            t.Errorf("Expected status 201, got %d", resp.StatusCode)
        }

        if err := json.NewDecoder(resp.Body).Decode(&createAnswerResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        answerID = createAnswerResp.Answers[0].ID
        if answerID == 0 {
            t.Error("Expected answer ID to be set")
        }

		// 11. Получение текущих ответов 
		resp, err = client.Get(fmt.Sprintf("%s/questions/%d", testServer.URL, questionID))
        if err != nil {
            t.Fatalf("Failed to get question with answers: %v", err)
        }
        defer resp.Body.Close()

        if err := json.NewDecoder(resp.Body).Decode(&getQuestionResp); err != nil {
            t.Fatalf("Failed to decode response: %v", err)
        }

        if len(getQuestionResp.Answers) == 0 {
            t.Errorf("Expected at least 1 answer, got %d", len(getQuestionResp.Answers))
        }

        // 11. Удаление вопроса
        req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/questions/%d", testServer.URL, questionID), nil)
        if err != nil {
            t.Fatalf("Failed to create delete request: %v", err)
        }

        resp, err = client.Do(req)
        if err != nil {
            t.Fatalf("Failed to delete question: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNoContent {
            t.Errorf("Expected status 204, got %d", resp.StatusCode)
        }

        // 12. Проверка, что вопрос удален
        resp, err = client.Get(fmt.Sprintf("%s/questions/%d", testServer.URL, questionID))
        if err != nil {
            t.Fatalf("Failed to get deleted question: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusNotFound {
            t.Errorf("Expected status 404 for deleted question, got %d", resp.StatusCode)
        }

		// 13. Проверка, что все ответы удалены
		for i := 0; i < len(getQuestionResp.Answers); i++ {

			resp, err = client.Get(fmt.Sprintf("%s/answers/%d", testServer.URL, getQuestionResp.Answers[i].ID))
			if err != nil {
				t.Fatalf("Failed to get answer: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("Expected status 404, got %d", resp.StatusCode)
			}
		}
    })
}
