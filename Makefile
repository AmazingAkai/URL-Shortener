.PHONY: build run migrate clean watch

APP_DIR=app

build:
	@echo "Building..."
	@cd $(APP_DIR) && go build -o build/main cmd/api/main.go

run:
	@cd $(APP_DIR) && go run cmd/api/main.go

migrate:
	@: $(or $(STEPS),$(error "Error: STEPS is required!"))
	@cd $(APP_DIR) && go run cmd/migrate/main.go -steps $(STEPS)

clean:
	@echo "Cleaning..."
	@cd $(APP_DIR) && rm -f build/main

watch:
	@command -v air >/dev/null 2>&1 || { echo >&2 "Error: air is not installed. Please install it first."; exit 1; }
	@cd $(APP_DIR) && air
