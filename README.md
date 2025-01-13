# Project github.com/AmazingAkai/URL-Shortener

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
docker compose up -d
```

Wait for docker to build the application and run it
