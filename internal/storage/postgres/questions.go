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
	res := models.QuestionWithAnswers{}
	if err := s.conn.First(&res, id).Error; err != nil {
        return res, err
    }
    err := s.conn.Model(&s.conn).Association("answers").Find(&res.Answers)
    return res, err
}

func(s *Storage) AllQuestions(ctx context.Context) ([]models.Question, error) {
	return gorm.G[models.Question](s.conn).Find(ctx)
}

func(s *Storage) DeleteQuestion(ctx context.Context, id int) (int, error) {
	return gorm.G[models.Question](s.conn).Where("id = ?", id).Delete(ctx)
}

func(s *Storage) Exist(ctx context.Context, id int) (bool, error) {
	return gorm.G[bool](s.conn).Where("id = ?", id).First(ctx)
}