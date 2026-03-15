package respond

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"hello": "world"}

	JSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got["hello"] != "world" {
		t.Errorf("expected hello=world, got hello=%s", got["hello"])
	}
}

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	OK(w, map[string]int{"count": 42})

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var got map[string]int
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got["count"] != 42 {
		t.Errorf("expected count=42, got count=%d", got["count"])
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, map[string]string{"id": "abc-123"})

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var got map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got["id"] != "abc-123" {
		t.Errorf("expected id=abc-123, got id=%s", got["id"])
	}
}

func TestNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	NoContent(w)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("expected empty body, got %q", w.Body.String())
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusNotFound, "resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var got struct {
		Error struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.Error.Status != http.StatusNotFound {
		t.Errorf("expected error status %d, got %d", http.StatusNotFound, got.Error.Status)
	}
	if got.Error.Message != "resource not found" {
		t.Errorf("expected message 'resource not found', got %q", got.Error.Message)
	}
}

func TestErrorWithDetails(t *testing.T) {
	w := httptest.NewRecorder()
	details := map[string]string{"field": "email", "reason": "invalid format"}
	ErrorWithDetails(w, http.StatusBadRequest, "validation failed", details)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var got struct {
		Error struct {
			Status  int               `json:"status"`
			Message string            `json:"message"`
			Details map[string]string `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.Error.Status != http.StatusBadRequest {
		t.Errorf("expected error status %d, got %d", http.StatusBadRequest, got.Error.Status)
	}
	if got.Error.Message != "validation failed" {
		t.Errorf("expected message 'validation failed', got %q", got.Error.Message)
	}
	if got.Error.Details["field"] != "email" {
		t.Errorf("expected details field=email, got %q", got.Error.Details["field"])
	}
	if got.Error.Details["reason"] != "invalid format" {
		t.Errorf("expected details reason='invalid format', got %q", got.Error.Details["reason"])
	}
}

func TestJSON_MarshalError(t *testing.T) {
	w := httptest.NewRecorder()
	// Channels cannot be marshalled to JSON.
	JSON(w, http.StatusOK, make(chan int))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
