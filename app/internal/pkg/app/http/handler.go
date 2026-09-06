package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type PromptClient interface {
	Prompt(ctx context.Context, prompt string) (string, error)
	Model() string
}

type Handler struct {
	ai PromptClient
}

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type PromptResponse struct {
	Answer string `json:"answer"`
	Model  string `json:"model"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Router interface {
	Get(pattern string, handlerFn http.HandlerFunc)
	Post(pattern string, handlerFn http.HandlerFunc)
}

func NewHandler(ai PromptClient) *Handler {
	return &Handler{ai: ai}
}

func (h *Handler) RegisterRoutes(r Router) {
	r.Get("/health", h.health)
	r.Post("/prompt", h.prompt)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) prompt(w http.ResponseWriter, r *http.Request) {
	var req PromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json body"})
		return
	}

	req.Prompt = strings.TrimSpace(req.Prompt)
	if req.Prompt == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "prompt is required"})
		return
	}

	answer, err := h.ai.Prompt(r.Context(), req.Prompt)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, ErrorResponse{Error: "llm request failed"})
		return
	}

	writeJSON(w, http.StatusOK, PromptResponse{
		Answer: answer,
		Model:  h.ai.Model(),
	})
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
