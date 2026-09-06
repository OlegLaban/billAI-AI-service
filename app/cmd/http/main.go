package main

import (
	"bytes"
	"errors"
	"log"
	"os"

	apphttp "github.com/OlegLaban/billAI-AI-service/app/internal/pkg/app/http"
	httpserver "github.com/OlegLaban/billAI-AI-service/app/pkg/http_server"
	"github.com/OlegLaban/billAI-AI-service/app/pkg/logger"
	ollama "github.com/OlegLaban/billAI-AI-service/app/pkg/ollama"
	"github.com/go-chi/chi/v5"
)

func main() {
	b, err := os.ReadFile("./.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	c := apphttp.ParseConfig(bytes.NewReader(b))
	l := logger.New()

	aiClient, err := ollama.New(ollama.Config{
		Host:   c.OllamaHost,
		Model:  c.OllamaModel,
		System: c.OllamaSystem,
	}.WithTemperature(c.OllamaTemperature))
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	apphttp.NewHandler(aiClient).RegisterRoutes(r)

	s := httpserver.NewServer(r, l)
	app := apphttp.NewApp(s, c)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
