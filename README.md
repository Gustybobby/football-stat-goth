# Project Setup
This project is mainly built with Golang and uses PostgreSQL as the database.

## Golang (sqlc, templ, air)
- Install sqlc (type-safe code generation from SQL)
- templ (go HTML templating language)
- air (go live reloading)

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/air-verse/air@latest
```

## Node & NPM (For Tailwind CSS)

### Install Node.js
We only need Node.js for compiling Tailwind CSS classes

### Install Dev Dependencies
Install Tailwind CSS and Tailwind Motion (Rombo) as dev dependencies

```bash
npm install -d tailwindcss
npm install -d tailwindcss-motion
```

## PostgreSQL

### Schema
Create tables, enums, and triggers with [schema.sql](schema.sql)

### Seed (Optional)
Optionally, you can seed the database with [seed.sql](seed.sql). The seed contains all matches data from 2024/25 Premier League week 1 - 9.

## Environment Variables
See [env example](.env.example)

# Dev Mode
Start the application in dev mode with air

```bash
air
```

Compile templ, sql, CSS classes using

```bash
npm run build
```
