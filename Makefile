BIN_DIR=bin
APP_NAME=avito-test-task

build:
	go build -o $(BIN_DIR)/app ./cmd/app/main.go
	go build -o $(BIN_DIR)/migrator ./cmd/migrator/main.go

run: build
	./$(BIN_DIR)/migrator --storage-path=./storage/sqlite/pr.db --migrations-path=./internal/migrations
	./$(BIN_DIR)/app

