.PHONY: run migrate build-templ build-tailwind build watch-templ watch-tailwind watch

run:
	@go run cmd/api/main.go

migrate:
	@: $(or $(STEPS),$(error "Error: STEPS is required!"))
	@go run cmd/migrate/main.go -steps $(STEPS)

build-templ:
	@templ generate

build-tailwind:
	@npx tailwindcss -i static/css/style.css -o static/css/tailwind.css --minify

build:
	@go build -o build/main cmd/api/main.go

watch-templ:
	@templ generate -watch -proxy http://localhost:8080

watch-tailwind:
	@npx tailwindcss -i static/css/style.css -o static/css/tailwind.css --watch

watch:
	@air
