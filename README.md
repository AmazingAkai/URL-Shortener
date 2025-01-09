# Project github.com/AmazingAkai/URL-Shortener

## Setting up database

```sql
CREATE ROLE url_shortener WITH LOGIN PASSWORD 'url_shortener';
CREATE DATABASE url_shortener OWNER url_shortener;
CREATE EXTENSION pg_trgm;
```

### Running Migrations

```bash
make migrate STEPS=1 # number of steps
```

## Getting Started

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Clean up binary from the last build:
```bash
make clean
```

Live reload the application:
```bash
make watch
```

