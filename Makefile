APP_NAME=qqbot
BUILD_DIR=build

.PHONY: all build run clean deps test

all: build

deps:
	go mod download
	go mod tidy

build: deps
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/bot/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

dev:
	go run cmd/bot/main.go

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run
