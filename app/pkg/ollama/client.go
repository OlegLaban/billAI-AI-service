package ai

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

// Client communicates with an Ollama server.
//
// Create with [New], verify connectivity with [Client.Ping], then use
// prompt, chat, or conversation methods. A single Client is safe for
// concurrent use by multiple goroutines (the underlying HTTP client is shared).
type Client struct {
	cfg    Config
	ollama *api.Client
}

// New creates a client for the given configuration.
//
// Returns an error when Model is empty or Host is set but not a valid URL.
// When Host is empty, the address is read from the OLLAMA_HOST environment variable.
func New(cfg Config) (*Client, error) {
	if cfg.Model == "" {
		return nil, fmt.Errorf("ai: model is required")
	}

	var ollamaClient *api.Client
	var err error
	if cfg.Host == "" {
		ollamaClient, err = api.ClientFromEnvironment()
	} else {
		u, parseErr := url.Parse(cfg.Host)
		if parseErr != nil {
			return nil, fmt.Errorf("ai: invalid host: %w", parseErr)
		}
		ollamaClient = api.NewClient(u, http.DefaultClient)
	}
	if err != nil {
		return nil, fmt.Errorf("ai: create ollama client: %w", err)
	}

	return &Client{cfg: cfg, ollama: ollamaClient}, nil
}

// Model returns the configured model name.
func (c *Client) Model() string {
	return c.cfg.Model
}

// Ping checks that the Ollama server is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.ollama.Heartbeat(ctx)
}
