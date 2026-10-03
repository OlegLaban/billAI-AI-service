build:
	docker compose build

compile:
	mkdir -p ./app/build
	cd ./app && go build -o ./build/http ./cmd/http

compile-q:
	mkdir -p ./app/build/
	cd ./app && go build -o ./build/queue ./cmd/queue

up:
	docker compose up -d

pull:
	docker exec -i ollama ollama pull gemma3:1b
