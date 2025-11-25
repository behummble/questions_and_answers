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
	panic(s.server.ListenAndServe())
}

func(s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func newMux(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/questions", s.createQuestion)
	mux.HandleFunc("GET /api/v1/questions", s.getAllQuestions)
	mux.HandleFunc("GET /api/v1//questions/{id}", s.getQuestion)
	mux.HandleFunc("DELETE /api/v1/questions/{id}", s.deleteQuestion)

	mux.HandleFunc("POST /api/v1/questions/{id}/answers", s.createAnswer)
	mux.HandleFunc("GET /api/v1/answers/{id}", s.getAnswer)
	mux.HandleFunc("DELETE /api/v1/answers/{id}", s.deleteAnswer)
	
	return mux
}

func executeRequestBody(request *http.Request, log *slog.Logger) ([]byte, error) {
	body, err := request.GetBody()
	if err != nil {
		log.Error(
			"ExecutionBodyError", 
			slog.String("component", "io/Read"),
			slog.Any("error", err),
		)
		return nil, errors.New("ExecutionBodyError")
	}
	data, err := io.ReadAll(body)
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