package models

import (
	"time"
)

type Question struct {
	ID int
	Text string
	CreatedAt time.Time
}

type CreateQuestionRequest struct {
	ID int
	Text string
	CreatedAt time.Time
}

type GetQuestionsResponse struct {
	Questions []Question
}

type GetQuestionResponse struct {
	Question Question
	Answers[]Answer
}