BIN_NAME=server-exe
DEBUG=true

run: 
	@echo "Начинаю запускать сервер через makefile"
	@go run ./cmd/api

dev:
	@DEBUG=$(DEBUG) go run ./cmd/api

.PHONY: build
build: 
	@echo "Сборка начата"
	go build -o $(BIN_NAME) ./cmd/api
	@echo "Сборка завершена"

run-after-build: build
	./$(BIN_NAME)

clean:
	rm $(BIN_NAME)