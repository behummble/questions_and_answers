package http

import (
	"context"
	"log/slog"
	"net/http"
	"fmt"

	"github.com/behummble/Questions-answers/internal/config"
	"github.com/behummble/Questions-answers/internal/models"
)

type Server struct {
	log *slog.Logger
	server *http.Server
	service Service
}

type Service interface {
	NewQuestion(ctx context.Context, data []byte) error
	Question(ctx context.Context, id int) (models.GetQuestionResponse, error)
	AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error)
	DeleteQuestion(ctx context.Context, id int) error
	NewAnswer(ctx context.Context, answer []byte, questionID int) error 
	Answer(ctx context.Context, id int) (models.Answer, error) 
	DeleteAnswer(ctx context.Context, id int) error
}

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.ServerConfig, service Service) *Server {
	server := &Server{
		log: log,
		service: service,
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	router := newRouter()
	router.register(server)
	srv.Handler = router
	server.server = srv
	
	return server
}

func(s *Server) Start() {
	s.server.ListenAndServe()
}

func(s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func newMux(s *Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/questions", s.questions)
	mux.HandleFunc("/api/v1/email_task", s.answers)
	return mux
}