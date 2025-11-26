package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)


func(s *Server) CreateAnswer(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	data, err := executeRequestBody(request, s.log)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}

	if len(data) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, "Empty body")
	}

	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}

	s.log.Info(fmt.Sprintf("Recive a request to create answer for question with id: %d", id))

	res, err := s.service.NewAnswer(ctx, data, id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusCreated)
	writer.Write(bytes)
}

func(s *Server) GetAnswer(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	s.log.Info(fmt.Sprintf("Recive request for get answer with id: %d", id))
	res, err := s.service.Answer(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err.Error())
	}
	bytes := prepareResponse(res, s.log)
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func(s *Server) DeleteAnswer(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 30 * time.Second)
	defer cancel()
	id, err := getID(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	s.log.Info(fmt.Sprintf("Recive request for delete answer with id: %d", id))
	err = s.service.DeleteAnswer(ctx, id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err.Error())
		return
	}
	
	writer.WriteHeader(http.StatusNoContent)
}
