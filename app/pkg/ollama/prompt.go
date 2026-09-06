package ai

import (
	"context"
	"strings"

	"github.com/ollama/ollama/api"
)

// Prompt sends a single prompt via the Ollama Generate API and returns the full response.
// Config.System is passed as the system instruction for this request.
func (c *Client) Prompt(ctx context.Context, prompt string) (string, error) {
	stream := false
	req := &api.GenerateRequest{
		Model:   c.cfg.Model,
		Prompt:  prompt,
		System:  c.cfg.System,
		Stream:  &stream,
		Options: c.cfg.Options,
	}

	var sb strings.Builder
	err := c.ollama.Generate(ctx, req, func(resp api.GenerateResponse) error {
		sb.WriteString(resp.Response)
		return nil
	})
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

// PromptStream streams the Generate API response, calling onChunk for each non-empty fragment.
// Returning a non-nil error from onChunk aborts generation and propagates that error.
func (c *Client) PromptStream(ctx context.Context, prompt string, onChunk func(chunk string) error) error {
	req := &api.GenerateRequest{
		Model:   c.cfg.Model,
		Prompt:  prompt,
		System:  c.cfg.System,
		Options: c.cfg.Options,
	}

	return c.ollama.Generate(ctx, req, func(resp api.GenerateResponse) error {
		if resp.Response == "" {
			return nil
		}
		return onChunk(resp.Response)
	})
}
