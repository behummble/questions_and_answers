package http

import (
	"encoding/json"
	"net/http"
	"testing"
	"bytes"
	"github.com/behummble/Questions-answers/internal/models"
)

func TestExecuteRequestBodyWithBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/questions", bytes.NewReader([]byte("test")))
    if err != nil {
        t.Fatal(err)
    }

	if req.Body == nil {
		t.Errorf("Excpect not empty body of request: got %v want not nil",
            req.Body)
	}
}

func TestExecuteRequestBodyWithoutBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/questions", nil)
    if err != nil {
        t.Fatal(err)
    }

	if req.Body != nil {
		t.Errorf("Excpect empty body of request: got %v want not nil",
            req.Body)
	}
}

func TestPrepareResponse(t *testing.T) {
	m := models.CreateQuestionResponse{}
	res, err := json.Marshal(m)
	if err != nil {
        t.Fatal(err)
    }

	if string(res) == "" {
		t.Errorf("Empty string in json marshalin: got %s want not empty",
            res)
	}
}
func TestGetIDWithoutID(t *testing.T) {
	req, err := http.NewRequest("POST", "/questions", bytes.NewReader([]byte("test")))
	if err != nil {
        t.Fatal(err)
    }
	idStr := req.PathValue("id")
	if idStr != "" {
		t.Errorf("Excpect empty path value: got %s want empty",
            idStr)
	}
}