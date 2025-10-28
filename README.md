# shop-go-api

Work in progress RESTful API for an online shop.

## Table of Contents

- [About the Project](#about-the-project)
- [Installation](#installation)
- [Docs](#docs)
- [Planed features](#planed-features)

---

## About the Project

API used to organize an online shop.  
It currently supports:

- JWT-based authentication
- Admin support for updating or fetching user data

---

## Installation

Follow these steps to run the API.

### Prerequisites

- Go (1.25+)
- PostgreSQL
- Git

### Steps

1. **Clone the repository**

   ```bash
   git clone https://github.com/Simone-Samardzhiev/shop-api-go
   cd shop-go-api
   ```

2. **Run the setup CLI to register the first admin**

   ```bash
   go run ./cmd/setup/main.go
   ```

   or using the Taskfile:

   ```bash
   task run-setup
   ```

3. **Add environment variables (either via `.env` file or export them directly)**

   #### `.env` example:
   ```ini
   ENVIRONMENT=development
   PORT=:8080
   DATABASE_URL=postgres://postgres:username@host:port/database
   DATABASE_MAX_OPEN_CONNECTIONS=10
   DATABASE_MAX_IDLE_CONNECTIONS=10
   JWT_SECRET=secret
   JWT_REFRESH_TOKEN_EXPIRE_TIME=10h
   JWT_ACCESS_TOKEN_EXPIRE_TIME=10m
   JWT_ISSUER=my-app
   JWT_AUDIENCE=my-app-audience
   ```

   #### Or export directly:
   ```bash
   export ENVIRONMENT=development
   export PORT=:8080
   export DATABASE_URL=postgres://postgres:username@host:port/database
   export DATABASE_MAX_OPEN_CONNECTIONS=10
   export DATABASE_MAX_IDLE_CONNECTIONS=10
   export JWT_SECRET=secret
   export JWT_REFRESH_TOKEN_EXPIRE_TIME=10h
   export JWT_ACCESS_TOKEN_EXPIRE_TIME=10m
   export JWT_ISSUER=my-app
   export JWT_AUDIENCE=my-app-audience
   ```

4. **Run database migrations**

   ```bash
   task migrate-up
   ```

5. **Run the API**

   ```bash
   go run ./cmd/http/main.go
   ```

   or with the Taskfile:

   ```bash
   task run-api
   ```

---

## Docs

- Swagger docs:
    - [JSON version](docs/swagger.json)
    - [YAML version](docs/swagger.yaml)

- Apidog documentation:  
  👉 [View it here](https://jm9m3ngpy4.apidog.io)

---

## Planed features

- Product and inventory management
- Order creation and checkout system
- API rate limiting

---

## Development Tools

- [Task](https://taskfile.dev) — for running project tasks easily
- [golang-migrate](https://github.com/golang-migrate/migrate) — for database migrations