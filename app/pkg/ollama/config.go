package ai

// Config holds settings for the AI client.
//
// Use [Config.WithTemperature] and [Config.WithOption] to set model parameters
// without mutating the original value — both methods return a copy.
type Config struct {
	// Model is the Ollama model name (e.g. "llama3.2", "gemma2").
	Model string

	// System is an optional system prompt applied to every request.
	System string

	// Host is the Ollama server URL (e.g. "http://localhost:11434").
	// If empty, OLLAMA_HOST env var is used.
	Host string

	// Options are model-specific parameters (temperature, top_p, etc.).
	Options map[string]any
}

// WithOption returns a copy of Config with the given option set.
func (c Config) WithOption(key string, value any) Config {
	opts := make(map[string]any, len(c.Options)+1)
	for k, v := range c.Options {
		opts[k] = v
	}
	opts[key] = value
	c.Options = opts
	return c
}

// WithTemperature returns a copy of Config with temperature set.
func (c Config) WithTemperature(t float64) Config {
	return c.WithOption("temperature", t)
}
