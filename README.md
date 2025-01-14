# Project github.com/AmazingAkai/URL-Shortener

A fun URL shortener built with Go, HTMX, Templ, TailwindCSS and PostgreSQL (GoTTH) stack. This url shortener also comes with browser extensions for chrome (or chrome based browsers) and firefox.

## Setting up database

## Without Docker

```sql
CREATE ROLE url_shortener WITH LOGIN PASSWORD 'url_shortener';
CREATE DATABASE url_shortener OWNER url_shortener;
CREATE EXTENSION pg_trgm;
```

## With Docker

```bash
docker compose up -d db
```

## Configuring Environment

```bash
cp .env.example .env
```

Fill `.env` with your database credentials.

### Running Migrations

```bash
make migrate STEPS=1 # number of steps
```

## Running (Development)

Build the application

```bash
make build-tailwind
make build-templ
make build
```

Run the application

```bash
./build/main
```

## Running (Production)

Build the application

```bash
docker build -t url-shortener .
```

Wait for docker to build the application and run it using:

```bash
docker compose up -d
```
