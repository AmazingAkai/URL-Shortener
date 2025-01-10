.PHONY: run migrate build-templ build-tailwind build watch-templ watch-tailwind watch

run:
	$(MAKE) build-templ && $(MAKE) build-tailwind && go run cmd/api/main.go

migrate:
	@: $(or $(STEPS),$(error "Error: STEPS is required!"))
	go run cmd/migrate/main.go -steps $(STEPS)

build-templ:
	templ generate

build-tailwind:
	npx tailwindcss -i static/css/style.css -o static/css/tailwind.css --minify

build:
	$(MAKE) build-templ && $(MAKE) build-tailwind && go build -o build/main cmd/api/main.go

watch-templ:
	templ generate -watch -proxy http://localhost:8080

watch-tailwind:
	npx tailwindcss -i static/css/style.css -o static/css/tailwind.css --watch

watch:
	$(MAKE) watch-templ & $(MAKE) watch-tailwind & air
