package postgres

import (
	"context"
	"gorm.io/gorm"
	"github.com/behummble/Questions-answers/internal/models"
)

func(s *Storage) CreateQuestion(ctx context.Context, data models.Question) error {

}

func(s *Storage) Question(ctx context.Context, id int) (models.GetQuestionResponse, error) {

}

func(s *Storage) AllQuestions(ctx context.Context) (models.GetQuestionsResponse, error) {

}

func(s *Storage) DeleteQuestion(ctx context.Context, id int) error {

}

func(s *Storage) Exist(ctx context.Context, id int) (bool, error) {

}