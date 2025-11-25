package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func(s *Server) createQuestion(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 5 * time.Second)
	defer cancel()
	data, err := executeRequestBody(request, s.log)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
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

func(s *Server) getAllQuestions(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 5 * time.Second)
	defer cancel()
	res, err := s.service.AllQuestions(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func(s *Server) getQuestion(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 5 * time.Second)
	defer cancel()
	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	res, err := s.service.Question(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func(s *Server) deleteQuestion(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 5 * time.Second)
	defer cancel()
	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	err = s.service.DeleteQuestion(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	
	writer.WriteHeader(http.StatusNoContent)
}
