package queue

import (
	"context"
	"fmt"
)

type Message struct {
	Prompt string
}

type Log interface {
	Error(msg string, err error)
}

type Client interface {
	Handle(ctx context.Context, callback func(context.Context, Message) error) error
}

type PromptClient interface {
	Prompt(ctx context.Context, prompt string) (string, error)
	Model() string
}

type App struct {
	c  Client
	ai PromptClient
	l  Log
}

func New(c Client, ai PromptClient, l Log) *App {
	return &App{c: c, ai: ai, l: l}
}

func (a *App) Run(ctx context.Context) error {
	fmt.Println("Reading from queue was started")
	a.c.Handle(ctx, func(ctx context.Context, m Message) error {

		answ, err := a.ai.Prompt(ctx, m.Prompt)
		if err != nil {
			a.l.Error("can`t send data to AI", err)
			return err
		}
		fmt.Println("Answer", answ)
		return nil
	})
	return nil
}
