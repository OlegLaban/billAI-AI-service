package http

import (
	"errors"

	internalerrors "github.com/OlegLaban/billAI-AI-service/app/internal/domain/internal_errors"
)

type Server interface {
	Serve(port int) error
}

type App struct {
	s Server
	c Config
}

func NewApp(s Server, c Config) *App {
	return &App{s: s, c: c}
}

func (a *App) Run() error {
	if err := a.s.Serve(a.c.Port); err != nil {
		return errors.Join(internalerrors.ErrCantStartApp, err)
	}
	return nil
}
