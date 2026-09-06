package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakePromptClient struct {
	answer string
	model  string
}

func (c fakePromptClient) Prompt(_ context.Context, prompt string) (string, error) {
	return c.answer, nil
}

func (c fakePromptClient) Model() string {
	return c.model
}

func TestPromptReturnsAnswer(t *testing.T) {
	handler := NewHandler(fakePromptClient{answer: "Hello", model: "test-model"})
	req := httptest.NewRequest(http.MethodPost, "/prompt", strings.NewReader(`{"prompt":"Hi"}`))
	rec := httptest.NewRecorder()

	handler.prompt(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"answer":"Hello","model":"test-model"}` + "\n"
	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestPromptRequiresPrompt(t *testing.T) {
	handler := NewHandler(fakePromptClient{})
	req := httptest.NewRequest(http.MethodPost, "/prompt", strings.NewReader(`{"prompt":"   "}`))
	rec := httptest.NewRecorder()

	handler.prompt(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
