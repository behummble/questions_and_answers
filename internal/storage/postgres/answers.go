package postgres

import (
	"context"
	"gorm.io/gorm"
	"github.com/behummble/Questions-answers/internal/models"
)

func(s *Storage) CreateAnswer(ctx context.Context, data []*models.Answer) error {
	return s.conn.WithContext(ctx).Create(data).Error
}

func(s *Storage) GetAnswer(ctx context.Context, id int) (models.Answer, error) {
	return gorm.G[models.Answer](s.conn).Where("id = ?", id).First(ctx)
}

func(s *Storage) DeleteAnswer(ctx context.Context, id int) (int, error) {
	return gorm.G[models.Answer](s.conn).Where("id = ?", id).Delete(ctx)
}
