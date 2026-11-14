MAIN_PATH = cmd/app/main.go
BIN_NAME = bin/tokigura
GOOSE_ENV ?= .env.dev

.PHONY: run build test lint \
        mg-create mg-up mg-down mg-status \
        dev-clean

run:
	GO_DEV=1 go run $(MAIN_PATH)

build:
	go build -ldflags="-s -w" -trimpath -o $(BIN_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BIN_NAME)"

test:
	go test -v -race ./...

mg-create:
	@if [ -z "$(name)" ]; then \
		echo "Specify migration name. Example: make mg-create name=init"; \
		exit 1; \
	fi
	goose -env=$(GOOSE_ENV) create $(name) sql

mg-up:
	goose -env=$(GOOSE_ENV) up

mg-down:
	goose -env=$(GOOSE_ENV) down

mg-status:
	goose -env=$(GOOSE_ENV) status

dev-clean:
	@echo "Not implemented"
