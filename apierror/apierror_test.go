package apierror_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo_api/apierror"
	"todo_api/middleware"
)

func TestWriteErrorFormatting(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := apierror.BadRequest("invalid payload", errors.New("missing field"))
		apierror.Write(w, r, err)
	})

	loggedHandler := middleware.Logger(handler)
	rec := httptest.NewRecorder()
	loggedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp apierror.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Error != "invalid payload" {
		t.Errorf("expected error message %q, got %q", "invalid payload", resp.Error)
	}
}
