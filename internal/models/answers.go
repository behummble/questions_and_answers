package models

import (
	"time"
)

type Answer struct {
	ID int
	QuestionID int
	UserID int
	Text string
	CreatedAt time.Time
}

type CreateAnswerRequest struct {
	ID int
	UserID int
	Text string
	CreatedAt time.Time
}

type GetAnswerResponse struct {
	Answer Answer
}