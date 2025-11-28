package models

import (
	"time"
)

type Answer struct {
	ID int
	QuestionID int
	UserID string
	Text string
	CreatedAt time.Time
}

type CreateAnswerRequest struct {
	UserID string
	Texts []string
}

type CreateAnswerResponse struct {
	Answers []*Answer
}

type GetAnswerResponse struct {
	Answer Answer
}