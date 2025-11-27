package http

import (
	"context"
	"log/slog"
	"net/http"
	"fmt"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/behummble/Questions-answers/internal/config"
	"github.com/behummble/Questions-answers/internal/models"
)

type Server struct {
	log *slog.Logger
	server *http.Server
	service Service
}

type Service interface {
	NewQuestion(ctx context.Context, data []byte) (models.CreateQuestionResponse, error)
	Question(ctx context.Context, id int) (models.GetQuestionResponse, error)
	AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error)
	DeleteQuestion(ctx context.Context, id int) (error)
	NewAnswer(ctx context.Context, answer []byte, questionID int) (models.CreateAnswerResponse, error)
	Answer(ctx context.Context, id int) (models.GetAnswerResponse, error) 
	DeleteAnswer(ctx context.Context, id int) (error)
}

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.ServerConfig, service Service) *Server {
	server := &Server{
		log: log,
		service: service,
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	mux := newMux(server)
	srv.Handler = mux
	server.server = srv
	
	return server
}

func(s *Server) Start() {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}

func(s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func newMux(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /questions", s.CreateQuestion)
	mux.HandleFunc("GET /questions", s.GetAllQuestions)
	mux.HandleFunc("GET /questions/{id}", s.GetQuestion)
	mux.HandleFunc("DELETE /questions/{id}", s.DeleteQuestion)

	mux.HandleFunc("POST /questions/{id}/answers", s.CreateAnswer)
	mux.HandleFunc("GET /answers/{id}", s.GetAnswer)
	mux.HandleFunc("DELETE /answers/{id}", s.DeleteAnswer)
	
	return mux
}

func executeRequestBody(request *http.Request, log *slog.Logger) ([]byte, error) {
	if request.Body == nil {
		log.Error(
			"ExecutionBodyError", 
			slog.String("component", "io/Read"),
			slog.Any("error", "Empty body"),
		)
		return nil, errors.New("ExecutionBodyError")
	}
	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error(
			"ReadingRequestBodyError", 
			slog.String("component", "io/Read"),
			slog.Any("error", err),
		)
		return nil, errors.New("ReadingRequestBodyError")
	}

	return data, nil
}

func prepareResponse[T any](m T, log *slog.Logger) []byte {
	res, err := json.Marshal(m)
	if err != nil {
		log.Error(
			"MarshalingJSONError", 
			slog.String("component", "json/marshalling"),
			slog.Any("error", err),
		)
	}

	return res
}

func getID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errors.New("ParameterNotFound")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}