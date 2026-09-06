build:
	docker compose build

compile:
	mkdir -p ./app/build
	cd ./app && go build -o ./build/http ./cmd/http

up:
	docker compose up -d

pull:
	docker exec -i ollama ollama pull mistral:7b
