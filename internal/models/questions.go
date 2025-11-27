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
	Text string
}

type CreateQuestionResponse struct {
	Question Question
}

type GetQuestionsResponse struct {
	Questions []Question
}

type GetQuestionResponse struct {
	Question Question
	Answers []Answer
}

type QuestionWithAnswers struct {
	Question Question
	Answers []Answer
}