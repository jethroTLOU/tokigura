MAIN_PATH = cmd/main.go
BIN_NAME = bin/tokigura
MIGRATIONS_DIR = internal/database/migrations
DEV_DB_PATH = '/Users/jethro/Library/Application Support/tokigura_dev/database.db'

.PHONY: run build mg-create mg-down dev-clean

run:
	GO_DEV=1 go run $(MAIN_PATH)

build:
	go build -o $(BIN_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BIN_NAME)"


mg-create:
	@if [ -z "$(name)" ]; then \
		echo "Specify migration name. Example: make mg-create name=init"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

mg-down:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DEV_DB_PATH) down

dev-clean:
	rm -f $(DEV_DB_PATH)
