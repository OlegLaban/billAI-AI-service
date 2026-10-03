package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/OlegLaban/billAI-AI-service/app/internal/pkg/app/queue"
	"github.com/OlegLaban/billAI-AI-service/app/pkg/logger"
	ollama "github.com/OlegLaban/billAI-AI-service/app/pkg/ollama"
	"github.com/OlegLaban/billAI-AI-service/app/pkg/redis"
)

func main() {
	b, err := os.ReadFile("./.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	c := queue.ParseConfig(bytes.NewReader(b))
	l := logger.New()

	aiClient, err := ollama.New(ollama.Config{
		Host:   c.OllamaHost,
		Model:  c.OllamaModel,
		System: c.OllamaSystem,
	}.WithTemperature(c.OllamaTemperature))
	if err != nil {
		log.Fatal(err)
	}

	r := redis.New(c, l)
	fmt.Println("queue password - %s", c.Password)
	q := queue.New(r, aiClient, l)
	ctx := context.Background()
	q.Run(ctx)
}
