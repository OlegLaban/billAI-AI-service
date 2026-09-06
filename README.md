# billAI-AI-service

HTTP-сервис для отправки prompt-а в LLM через Ollama и возврата сгенерированного ответа.

## Требования

- Docker и Docker Compose
- Go 1.26+
- Запущенный Docker daemon

## Быстрый Запуск

Поднять Ollama и Swagger UI:

```bash
make up
```

Скачать модель, которая используется сервисом по умолчанию:

```bash
make pull
```

Собрать Go-приложение:

```bash
make compile
```

Запустить HTTP-сервис:

```bash
cd app
./build/http
```

По умолчанию сервис стартует на `http://localhost:8080`.

## Конфигурация

Конфиг лежит в `app/.env`:

```env
OLLAMA_HOST=http://127.0.0.1:11434
OLLAMA_MODEL=qwen2.5-coder:1.5b
```

Дополнительно можно задать:

```env
PORT=8080
OLLAMA_SYSTEM_PROMPT=Отвечай кратко и по делу.
OLLAMA_TEMPERATURE=0.7
```

После изменения `app/.env` нужно перезапустить HTTP-сервис.

## Использование API

Health check:

```bash
curl http://localhost:8080/health
```

Запрос к LLM:

```bash
curl -X POST http://localhost:8080/prompt \
  -H "Content-Type: application/json" \
  -d '{"prompt":"Сколько планет в Солнечной системе?"}'
```

Пример успешного ответа:

```json
{
  "answer": "В Солнечной системе 8 планет.",
  "model": "qwen2.5-coder:1.5b"
}
```

Если тело запроса некорректное или `prompt` пустой, сервис вернёт ошибку:

```json
{
  "error": "prompt is required"
}
```

## Swagger UI

Swagger UI поднимается вместе с compose:

```bash
docker compose up -d swagger-ui
```

Открыть документацию:

```text
http://localhost:8081
```

OpenAPI-спецификация находится в `docs/openapi.yaml`.

## Полезные Команды

```bash
make up       # поднять Docker-сервисы
make pull     # скачать модель в Ollama
make compile  # собрать бинарник ./app/build/http
```
