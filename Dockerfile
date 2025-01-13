FROM golang:alpine AS fetch-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download


FROM node:latest AS tailwind-stage
WORKDIR /app
COPY --chown=65532:65532 . .
RUN npm install && \
    npx tailwindcss -i static/css/style.css -o static/css/tailwind.css --minify

FROM ghcr.io/a-h/templ:latest AS generate-stage
WORKDIR /app
COPY --from=tailwind-stage /app ./
RUN ["templ", "generate"]

FROM golang:alpine AS build-stage
WORKDIR /app
COPY --from=generate-stage /app ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o /bin/entrypoint cmd/api/main.go

FROM alpine:latest AS deploy-stage
WORKDIR /app
COPY --from=build-stage /bin/entrypoint /bin/entrypoint
COPY --from=build-stage /app/static /app/static
EXPOSE 8080
ENTRYPOINT ["entrypoint"]
