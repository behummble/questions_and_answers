package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func(s *Server) CreateQuestion(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	data, err := executeRequestBody(request, s.log)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}

	s.log.Info("Recive request to create question")

	if len(data) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Empty body")
	}

	res, err := s.service.NewQuestion(ctx, data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusCreated)
	writer.Write(bytes)
}

func(s *Server) GetAllQuestions(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	res, err := s.service.AllQuestions(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	s.log.Info("Recive request to get all question")
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func(s *Server) GetQuestion(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	s.log.Info(fmt.Sprintf("Recive a request to get question with id: %d", id))
	res, err := s.service.Question(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func(s *Server) DeleteQuestion(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	s.log.Info(fmt.Sprintf("Recive a request to delete question with id: %d", id))
	err = s.service.DeleteQuestion(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	
	writer.WriteHeader(http.StatusNoContent)
}
