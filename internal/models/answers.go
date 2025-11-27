package models

import (
	"time"
)

type Answer struct {
	ID int
	Question `gorm:"foreignKey:question_id"`
	QuestionID int
	UserID string
	Text string
	CreatedAt time.Time
}

type CreateAnswerRequest struct {
	UserID string
	Text string
}

type CreateAnswerResponse struct {
	Answer Answer
}

type GetAnswerResponse struct {
	Answer Answer
}