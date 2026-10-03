package queue

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Domain   string `env:"QUEUE_HOST"`
	Port     int    `env:"QUEUE_PORT"`
	Password string `env:"QUEUE_PASS"`
	DBIndex  int    `env:"QUEUE_DB"`

	OllamaHost        string  `env:"OLLAMA_HOST"`
	OllamaModel       string  `env:"OLLAMA_MODEL" envDefault:"mistral:7b"`
	OllamaSystem      string  `env:"OLLAMA_SYSTEM_PROMPT"`
	OllamaTemperature float64 `env:"OLLAMA_TEMPERATURE" envDefault:"0.7"`
}

func ParseConfig(r io.Reader) Config {
	loadEnv(r)

	var c Config
	err := env.Parse(&c)
	if err != nil {
		panic("Can`t read env file")
	}

	return c
}

func loadEnv(r io.Reader) {
	if r == nil {
		return
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		_ = os.Setenv(key, value)
	}
}
