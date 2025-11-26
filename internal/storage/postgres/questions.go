package postgres

import (
	"context"
	"gorm.io/gorm"
	"github.com/behummble/Questions-answers/internal/models"
)

func(s *Storage) CreateQuestion(ctx context.Context, data *models.Question) error {
	return gorm.G[models.Question](s.conn).Create(ctx, data)
}

func(s *Storage) Question(ctx context.Context, id int) (models.QuestionWithAnswers, error) {
	return gorm.G[models.QuestionWithAnswers](s.conn).Where("id = ?", id).First(ctx)
}

func(s *Storage) AllQuestions(ctx context.Context) ([]models.Question, error) {
	return gorm.G[models.Question](s.conn).Find(ctx)
}

func(s *Storage) DeleteQuestion(ctx context.Context, id int) error {
	_, err := gorm.G[models.Answer](s.conn).Where("id = ?", id).Delete(ctx)
	return err
}

func(s *Storage) Exist(ctx context.Context, id int) (bool, error) {
	return gorm.G[bool](s.conn).Where("id = ?", id).First(ctx)
}