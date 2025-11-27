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
	question, err := gorm.G[models.Question](s.conn).Where("id = ?", id).First(ctx)
	if err != nil {
		return models.QuestionWithAnswers{}, err
	}
	answers, err := gorm.G[models.Answer](s.conn).Where("question_id = ?", id).Find(ctx)

	res := models.QuestionWithAnswers{
		Question: question,
		Answers: answers,
	}

	return res, err
}

func(s *Storage) AllQuestions(ctx context.Context) ([]models.Question, error) {
	return gorm.G[models.Question](s.conn).Find(ctx)
}

func(s *Storage) DeleteQuestion(ctx context.Context, id int) (int, error) {
	return gorm.G[models.Question](s.conn).Where("id = ?", id).Delete(ctx)
}

func(s *Storage) Exist(ctx context.Context, id int) (models.Question, error) {
	return gorm.G[models.Question](s.conn).Where("id = ?", id).First(ctx)
}