# Project github.com/AmazingAkai/URL-Shortener

## Setting up database

```sql
CREATE ROLE url_shortener WITH LOGIN PASSWORD 'url_shortener';
CREATE DATABASE url_shortener OWNER url_shortener;
CREATE EXTENSION pg_trgm;
```

## Getting Started

Enter the app directory
```bash
cd app
```

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

