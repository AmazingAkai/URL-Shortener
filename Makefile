.PHONY: build run migrate clean watch


build:
	@echo "Building..."
	@go build -o build/main cmd/api/main.go

run:
	@go run cmd/api/main.go

migrate:
	@: $(or $(STEPS),$(error "Error: STEPS is required!"))
	@go run cmd/migrate/main.go -steps $(STEPS)

clean:
	@echo "Cleaning..."
	@rm -f build/main

watch:
	@command -v air >/dev/null 2>&1 || { echo >&2 "Error: air is not installed. Please install it first."; exit 1; }
	@air
