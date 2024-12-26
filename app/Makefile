.PHONY: build run migrate clean watch

build:
	@echo "Building..."
	@go build -o build/main cmd/api/main.go

run:
	@go run cmd/api/main.go

migrate:
	@if [ -z "$(STEPS)" ]; then \
		echo "Error: STEPS is required!"; \
		exit 1; \
	fi
	@go run cmd/migrate/main.go -steps $(STEPS)

clean:
	@echo "Cleaning..."
	@rm -f build/main

watch:
	@air


