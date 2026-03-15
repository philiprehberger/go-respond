package respond

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProblem(t *testing.T) {
	w := httptest.NewRecorder()
	Problem(w, http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var got ProblemDetails
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if got.Status != http.StatusBadRequest {
		t.Errorf("expected problem status %d, got %d", http.StatusBadRequest, got.Status)
	}
}

func TestProblemWithOptions(t *testing.T) {
	w := httptest.NewRecorder()
	Problem(w, http.StatusForbidden,
		WithType("https://example.com/problems/forbidden"),
		WithTitle("Forbidden"),
		WithDetail("You do not have access to this resource"),
		WithInstance("/accounts/12345"),
		WithExtension("account_id", "12345"),
	)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if got["type"] != "https://example.com/problems/forbidden" {
		t.Errorf("expected type URI, got %v", got["type"])
	}
	if got["title"] != "Forbidden" {
		t.Errorf("expected title 'Forbidden', got %v", got["title"])
	}
	if int(got["status"].(float64)) != http.StatusForbidden {
		t.Errorf("expected status %d, got %v", http.StatusForbidden, got["status"])
	}
	if got["detail"] != "You do not have access to this resource" {
		t.Errorf("expected detail message, got %v", got["detail"])
	}
	if got["instance"] != "/accounts/12345" {
		t.Errorf("expected instance URI, got %v", got["instance"])
	}
	if got["account_id"] != "12345" {
		t.Errorf("expected extension account_id=12345, got %v", got["account_id"])
	}
}

func TestProblemContentType(t *testing.T) {
	w := httptest.NewRecorder()
	Problem(w, http.StatusInternalServerError)

	ct := w.Header().Get("Content-Type")
	if ct != "application/problem+json" {
		t.Errorf("expected Content-Type application/problem+json, got %s", ct)
	}
}
