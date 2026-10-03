package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/OlegLaban/billAI-AI-service/app/internal/pkg/app/queue"
	client "github.com/go-redis/redis"
)

type Logger interface {
	Error(msg string, err error)
}

type RedisClient struct {
	c *client.Client
	l Logger
}

func New(config queue.Config, l Logger) *RedisClient {
	c := client.NewClient(&client.Options{
		Addr: fmt.Sprintf("%s:%d", config.Domain, config.Port),
		DB:   config.DBIndex,
	})
	return &RedisClient{c: c, l: l}
}

func (rc *RedisClient) Handle(ctx context.Context, callback func(context.Context, queue.Message) error) error {
	for {
		result, err := rc.c.BLPop(0*time.Second, "bill_ai_queue").Result()
		if err != nil {
			rc.l.Error("can`t read data from queue", err)
			return err
		}
		msg := queue.Message{}
		err = json.NewDecoder(strings.NewReader(result[1])).Decode(&msg)
		if err != nil {
			rc.l.Error("can`t decode message from queue", err)
			return err
		}
		callback(ctx, msg)
	}

}
