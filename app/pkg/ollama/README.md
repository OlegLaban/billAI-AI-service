# Go SDK для Ollama

Go SDK для работы с [Ollama](https://ollama.com): одиночные промпты, чат, стриминг и многоходовые диалоги с сохранением истории.

Пакет расположен в `pkg/ollama`, но объявлен как `package ai`. Рекомендуемый импорт:

```go
import ollama "github.com/OlegLaban/billAI-AI-service/app/pkg/ollama"
```

Без алиаса идентификатор пакета в коде будет `ai`, а не `ollama`.

Внутри используется официальный клиент [`github.com/ollama/ollama/api`](https://github.com/ollama/ollama).

## Содержание

- [Требования](#требования)
- [Быстрый старт](#быстрый-старт)
- [Архитектура](#архитектура)
- [Конфигурация](#конфигурация)
- [Справочник API](#справочник-api)
- [Примеры](#примеры)
- [Сравнение режимов](#сравнение-режимов)
- [Обработка ошибок](#обработка-ошибок)
- [Docker](#docker)
- [Тесты и godoc](#тесты-и-godoc)

## Требования

- Go 1.26+
- Запущенный сервер Ollama (локально или в Docker)
- Скачанная модель (`ollama pull llama3.2` или `make pull` для `mistral:7b`)

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"
    "log"

    ollama "github.com/OlegLaban/billAI-AI-service/app/pkg/ollama"
)

func main() {
    client, err := ollama.New(ollama.Config{
        Model:  "llama3.2",
        System: "Отвечай кратко и по делу.",
    }.WithTemperature(0.7))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    if err := client.Ping(ctx); err != nil {
        log.Fatalf("ollama недоступен: %v", err)
    }

    answer, err := client.Prompt(ctx, "Сколько планет в Солнечной системе?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(answer)
}
```

## Архитектура

```
┌─────────────────────────────────────────────────────────┐
│                     Ваше приложение                      │
└──────────────────────────┬──────────────────────────────┘
                           │
              ┌────────────▼────────────┐
              │   pkg/ollama (ai)       │
              │                         │
              │  Config ──► Client      │
              │              │          │
              │    ┌─────────┼─────────┐│
              │    ▼         ▼         ▼│
              │ Prompt   Chat    Conversation
              │ Stream   Stream       Send
              └────────────┬────────────┘
                           │
              ┌────────────▼────────────┐
              │ github.com/ollama/ollama│
              │         /api            │
              └────────────┬────────────┘
                           │ HTTP
              ┌────────────▼────────────┐
              │   Ollama server :11434  │
              └─────────────────────────┘
```

SDK предоставляет три уровня абстракции:

| Уровень | Назначение |
|---------|------------|
| **Generate** (`Prompt`, `PromptStream`) | Один запрос — один ответ, без истории |
| **Chat** (`Chat`, `ChatStream`) | Список сообщений передаёт вызывающий код |
| **Conversation** | История хранится внутри объекта между вызовами `Send` |

## Конфигурация

### Config

| Поле | Тип | Обязательное | Описание |
|------|-----|:------------:|----------|
| `Model` | `string` | да | Имя модели Ollama (`llama3.2`, `mistral:7b` и т.д.) |
| `System` | `string` | нет | Системный промпт для всех запросов |
| `Host` | `string` | нет | URL сервера. Если пусто — берётся из `OLLAMA_HOST` |
| `Options` | `map[string]any` | нет | Параметры модели (`temperature`, `top_p`, `num_ctx` и др.) |

### Методы конфигурации

```go
cfg := ollama.Config{Model: "llama3.2"}

cfg = cfg.WithTemperature(0.7)
cfg = cfg.WithOption("top_p", 0.9)
cfg = cfg.WithOption("num_ctx", 4096)
```

`WithOption` и `WithTemperature` возвращают **копию** конфига — исходный объект не меняется.

### Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `OLLAMA_HOST` | `http://127.0.0.1:11434` | Адрес сервера (если `Config.Host` пуст) |

Явный `Host` в конфиге имеет приоритет:

```go
client, _ := ollama.New(ollama.Config{
    Model: "llama3.2",
    Host:  "http://ollama:11434", // Docker-сеть app-net
})
```

## Справочник API

### Типы

| Имя | Описание |
|-----|----------|
| `Config` | Настройки клиента: модель, system prompt, хост, options |
| `Client` | Клиент Ollama; создаётся через `New` |
| `Conversation` | Сессия с внутренней историей сообщений |
| `Message` | Одно сообщение чата (`Role` + `Content`) |
| `Role` | Роль автора: `RoleSystem`, `RoleUser`, `RoleAssistant` |

### Функции и методы

| Сигнатура | Описание |
|-----------|----------|
| `New(cfg Config) (*Client, error)` | Создаёт клиент. Ошибка, если `Model` пуст или `Host` невалиден |
| `(Config) WithOption(key, value) Config` | Копия конфига с дополнительным параметром модели |
| `(Config) WithTemperature(t float64) Config` | Копия конфига с заданной temperature |
| `(c *Client) Model() string` | Имя модели из конфигурации |
| `(c *Client) Ping(ctx) error` | Проверка доступности сервера (`Heartbeat`) |
| `(c *Client) Prompt(ctx, prompt) (string, error)` | Generate API: один промпт, полный ответ |
| `(c *Client) PromptStream(ctx, prompt, onChunk) error` | Generate API: стриминг по фрагментам |
| `(c *Client) Chat(ctx, messages) (string, error)` | Chat API: список сообщений → ответ ассистента |
| `(c *Client) ChatStream(ctx, messages, onChunk) error` | Chat API: стриминг ответа |
| `(c *Client) NewConversation() *Conversation` | Новая сессия с опциональным system prompt |
| `(conv *Conversation) Send(ctx, content) (string, error)` | Отправить сообщение, получить ответ, сохранить историю |
| `(conv *Conversation) History() []Message` | Копия текущей истории |
| `(conv *Conversation) Reset()` | Очистить историю (system prompt сохраняется) |

## Примеры

### Стриминг промпта

```go
err := client.PromptStream(ctx, "Напиши короткое стихотворение про Go", func(chunk string) error {
    fmt.Print(chunk)
    return nil
})
```

Ошибка из `onChunk` прерывает генерацию.

### Чат с явной историей

```go
messages := []ollama.Message{
    {Role: ollama.RoleUser, Content: "Привет!"},
    {Role: ollama.RoleAssistant, Content: "Здравствуйте!"},
    {Role: ollama.RoleUser, Content: "Расскажи про goroutines"},
}

reply, err := client.Chat(ctx, messages)
```

`Config.System` автоматически добавляется первым system-сообщением.

### Диалог с памятью

```go
conv := client.NewConversation()

reply, err := conv.Send(ctx, "Назови три необычных животных")
// ...
reply, err = conv.Send(ctx, "Какое из них самое опасное?")

// Посмотреть историю
for _, msg := range conv.History() {
    fmt.Printf("%s: %s\n", msg.Role, msg.Content)
}

conv.Reset() // начать заново, system prompt останется
```

При ошибке запроса последнее user-сообщение **откатывается** из истории.

### Таймаут запроса

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

answer, err := client.Prompt(ctx, prompt)
```

## Сравнение режимов

| Режим | API Ollama | История | System prompt |
|-------|------------|---------|---------------|
| `Prompt` / `PromptStream` | Generate | нет | через `Config.System` |
| `Chat` / `ChatStream` | Chat | передаёт вызывающий | добавляется автоматически |
| `Conversation.Send` | Chat | хранится внутри | добавляется при создании |

**Когда что использовать:**

- **Prompt** — классификация, суммаризация, разовые вопросы без контекста
- **Chat** — когда история управляется снаружи (БД, сессия HTTP)
- **Conversation** — интерактивный CLI или прототип чат-бота

## Обработка ошибок

Все методы запросов принимают `context.Context`. Типичные причины ошибок:

| Ситуация | Что делать |
|----------|------------|
| Сервер недоступен | Проверить `Ping`, `OLLAMA_HOST`, Docker |
| Модель не найдена | `ollama pull <model>` |
| Таймаут / отмена контекста | Увеличить timeout или проверить cancel |
| Ошибка в `onChunk` | Стриминг прерывается, ошибка возвращается вызывающему |

```go
answer, err := client.Prompt(ctx, prompt)
if err != nil {
    log.Printf("запрос не удался: %v", err)
}
```

## Docker

В корне репозитория:

```bash
make up      # поднять контейнер ollama
make pull    # скачать mistral:7b
```

Контейнер слушает порт `11434` в сети `app-net`. С хоста:

```bash
export OLLAMA_HOST=http://127.0.0.1:11434
```

Из другого контейнера в той же сети:

```go
Host: "http://ollama:11434"
```

## Тесты и godoc

```bash
cd app

# Unit-тесты
go test ./pkg/ollama/...

# Документация пакета в терминале
go doc github.com/OlegLaban/billAI-AI-service/app/pkg/ollama

# Локальный godoc-сервер
go doc -http=:6060
# открыть http://localhost:6060/pkg/github.com/OlegLaban/billAI-AI-service/app/pkg/ollama/
```

Тесты покрывают валидацию конфигурации и примеры (`Example*` в `examples_test.go`). Интеграционных тестов с живым Ollama нет.
