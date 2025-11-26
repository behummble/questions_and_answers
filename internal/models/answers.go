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
	UserID int
	Text string
}

type CreateAnswerResponse struct {
	Answer Answer
}

type GetAnswerResponse struct {
	Answer Answer
}