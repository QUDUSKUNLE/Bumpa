# Bumpa API

Backend API for the Bumpa application, built with Go, Echo, PostgreSQL, SQLC, Paystack, and an event-driven outbox architecture.

## Tech Stack

* **Go**
* **Echo** — HTTP framework
* **PostgreSQL** — primary database
* **pgx** — PostgreSQL driver
* **SQLC** — type-safe SQL code generation
* **golang-migrate** — database migrations
* **Paystack** — payment and transfer processing
* **Docker** — containerization
* **Make** — development and migration commands

## Project Structure

```text
.
├── adapters/
│   ├── config/
│   ├── db/
│   │   ├── migrations/
│   │   └── queries/
│   ├── events/
│   ├── handlers/
│   ├── payments/
│   ├── repositories/
│   └── routes/
│
├── cmd/
│   └── seed/                 # Database seed command
│
├── core/
│   ├── domain/
│   ├── services/
│   │   ├── badges/
│   │   ├── cashback/
│   │   └── outboxprocessor/
│   └── utils/
│
├── bin/
├── .env
├── .env.sample
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── docker-compose.migrate.yml
├── sqlc.json
└── go.mod
```

## Requirements

* Go
* Docker
* Docker Compose
* PostgreSQL
* Make

## Environment Variables

Create your local environment file:

```bash
cp .env.sample .env
```

Example `.env.sample`:

```env
DATABASE_URL=postgres://username:password@localhost:5432/bumpa?sslmode=disable

HTTP_PORT=8081

PAYSTACK_WEBHOOK_SECRET=

CASHBACK_AMOUNT_KOBO=30000

# PAYSTACK
PAYSTACK_SECRET_KEY=sk_test_your_secret_key
PAYSTACK_PUBLIC_KEY=pk_test_your_public_key
PAYSTACK_BASE_URL=https://api.paystack.co
```

### Environment Variables

| Variable                  | Description                               |
| ------------------------- | ----------------------------------------- |
| `DATABASE_URL`            | PostgreSQL connection string              |
| `HTTP_PORT`               | Port on which the API listens             |
| `PAYSTACK_WEBHOOK_SECRET` | Secret used to validate Paystack webhooks |
| `CASHBACK_AMOUNT_KOBO`    | Cashback amount in kobo                   |
| `PAYSTACK_SECRET_KEY`     | Paystack secret API key                   |
| `PAYSTACK_PUBLIC_KEY`     | Paystack public API key                   |
| `PAYSTACK_BASE_URL`       | Paystack API base URL                     |

> **Security:** Never commit `.env` or real Paystack credentials to Git.

Ensure `.env` is included in `.gitignore`:

```gitignore
.env
```

## Installation

Clone the repository:

```bash
git clone <repository-url>
cd Bumpa
```

Install project dependencies and development tools:

```bash
make setup
```

This installs:

* SQLC `v1.27.0`
* golang-migrate `v4.17.0`

into the local `bin/` directory.

## Database Setup

Configure PostgreSQL and set the connection string in `.env`:

```env
DATABASE_URL=postgres://username:password@localhost:5432/bumpa?sslmode=disable
```

Run all pending migrations:

```bash
make migrate-up
```

Check the current migration version:

```bash
make migration-version
```

## SQLC

Generate the SQLC code:

```bash
make generate
```

Run this whenever SQL queries or database definitions used by SQLC are changed.

## Database Seeding

The project provides a dedicated seed command for creating development/sample data.

Run:

```bash
make seed
```

The seed command creates the sample user account if it does not already exist.

This is intentionally separate from application startup so that starting the API does not repeatedly attempt to insert seed data.

## Running the Application

After configuring the environment and running migrations:

```bash
go run .
```

The API listens on:

```text
http://localhost:8081
```

The port can be changed through:

```env
HTTP_PORT=8081
```

## API Endpoints

### Public Endpoints

| Method | Endpoint                    | Description                           |
| ------ | --------------------------- | ------------------------------------- |
| `GET`  | `/`                         | Returns the API home/welcome response |
| `GET`  | `/health`                   | Health check                          |
| `GET`  | `/users/:user/achievements` | Returns a user's achievements         |
| `POST` | `/users/purchases`          | Creates a purchase for a user         |

### `GET /`

Returns the API home response.

```http
GET /
```

### `GET /health`

Checks whether the API is running.

```http
GET /health
```

Example:

```bash
curl http://localhost:8081/health
```

### `GET /users/:user/achievements`

Returns the user's unlocked achievements and the next available achievement for each achievement group.

```http
GET /users/{user}/achievements
```

Example response:

```json
{
  "unlocked_achievements": [
    "First Purchase",
    "Three Purchases",
    "Five Purchases"
  ],
  "next_available_achievements": [
    "Ten Purchases"
  ]
}
```

If the user has unlocked every achievement in a group, that group does not contribute an item to `next_available_achievements`.

### `POST /users/purchases`

Creates a purchase for a user.

```http
POST /users/purchases
```

A successful purchase may trigger the achievement and cashback workflow.

## Achievement System

Achievements are defined using:

```sql
CREATE TABLE achievements (
  code TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  achievement_group TEXT NOT NULL,
  position INT NOT NULL,
  UNIQUE (achievement_group, position)
);
```

User achievements are stored using:

```sql
CREATE TABLE user_achievements (
  user_id UUID NOT NULL REFERENCES users(id),
  achievement_code TEXT NOT NULL REFERENCES achievements(code),
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, achievement_code)
);
```

The `achievement_group` and `position` fields determine achievement progression.

Example:

| Code               | Name             | Group    | Position |
| ------------------ | ---------------- | -------- | -------: |
| `first_purchase`   | First Purchase   | shopping |        1 |
| `three_purchases`  | Three Purchases  | shopping |        2 |
| `five_purchases`   | Five Purchases   | shopping |        3 |
| `ten_purchases`    | Ten Purchases    | shopping |        4 |
| `twenty_purchases` | Twenty Purchases | shopping |        5 |

## Cashback Flow

Cashback is triggered through the event-driven architecture.

```text
User Purchase
     │
     ▼
Achievement Evaluation
     │
     ▼
AchievementUnlocked
     │
     ▼
Outbox Event
     │
     ▼
Outbox Processor
     │
     ▼
Event Bus
     │
     ▼
BadgeUnlocked
     │
     ▼
Cashback Service
     │
     ▼
Paystack Transfer
```

The application subscribes to:

```text
AchievementUnlocked
BadgeUnlocked
```

When a badge is unlocked, the cashback service attempts to send the configured cashback amount through Paystack.

## Outbox Processor

The outbox processor starts with the application:

```go
outboxProcessor := outboxprocessor.NewOutboxProcessor(repo, bus)

go outboxProcessor.Run(ctx)
```

The outbox pattern allows domain events to be persisted and processed asynchronously.

Failed events can be retried by the processor.

## Paystack

The application uses Paystack for cashback transfers.

Configure:

```env
PAYSTACK_SECRET_KEY=sk_test_your_secret_key
PAYSTACK_PUBLIC_KEY=pk_test_your_public_key
PAYSTACK_BASE_URL=https://api.paystack.co
```

Cashback amount is configured in kobo:

```env
CASHBACK_AMOUNT_KOBO=30000
```

`30000` kobo is equivalent to **₦300**.

### Paystack Transfer Recipient

Paystack transfers require a valid **recipient code**.

For example:

```text
RCP_xxxxxxxxxxxxx
```

The recipient code is not the same as:

* User ID
* Bank account number
* Phone number
* Email address

The recipient code should be stored against the user's payment account and supplied to Paystack as the transfer `recipient`.

If Paystack returns:

```text
invalid_transfer_recipient
```

verify that the value being sent as `recipient` is a valid Paystack recipient code.

## Docker

Build the application image:

```bash
docker build -t bumpa-api .
```

Run the container:

```bash
docker run --rm \
  --env-file .env \
  -p 8081:8081 \
  bumpa-api
```

### PostgreSQL and Docker

If PostgreSQL runs on the host machine, `localhost` inside the API container refers to the container itself.

If PostgreSQL is another Docker Compose service, use its service name.

For example:

```env
DATABASE_URL=postgres://username:password@postgres:5432/bumpa?sslmode=disable
```

## **Testing**

The project uses Go's built-in testing framework and `testify` for assertions.

### **Run All Tests**

Run the complete test suite:

```bash
go test ./...
```

## Docker Compose

Start the application:

```bash
docker compose up --build
```

Run in the background:

```bash
docker compose up --build -d
```

Stop the services:

```bash
docker compose down
```

## Docker Migrations

Run migrations using Docker:

```bash
make migrate-docker-up
```

Check migration version:

```bash
make migrate-docker-version
```

Rollback:

```bash
make migrate-docker-down
```

Force a migration version:

```bash
make migrate-docker-force V=3
```

## Production Migrations

Production uses:

```env
DATABASE_URL=postgres://...
```

Run:

```bash
make migrate-prod-up
```

Check the production migration version:

```bash
make migrate-prod-version
```

> **Warning:** Always verify `PROD_DATABASE_URL` before running production migrations.

## Migration Commands

### Run migrations

```bash
make migrate-up
```

### Rollback one migration

```bash
make migrate-down
```

### Check migration version

```bash
make migration-version
```

### Force migration version

```bash
make force-migrate V=3
```

### Roll back to a specific version

```bash
make migrate-rollback V=2
```

### Create a migration

```bash
make create-migration NAME=add_user_table
```

## Make Commands

| Command                               | Description                                   |
| ------------------------------------- | --------------------------------------------- |
| `make setup`                          | Install Go dependencies and development tools |
| `make generate`                       | Generate SQLC code                            |
| `make seed`                           | Seed development/sample data                  |
| `make migrate-up`                     | Run pending local migrations                  |
| `make migrate-down`                   | Rollback one local migration                  |
| `make migration-version`              | Show current migration version                |
| `make force-migrate V=VERSION`        | Force migration version                       |
| `make migrate-rollback V=VERSION`     | Roll back to a specific version               |
| `make migrate-docker-up`              | Run migrations using Docker                   |
| `make migrate-docker-down`            | Rollback Docker migration                     |
| `make migrate-docker-force V=VERSION` | Force Docker migration version                |
| `make migrate-docker-version`         | Show Docker migration version                 |
| `make migrate-prod-up`                | Run production migrations                     |
| `make migrate-prod-version`           | Show production migration version             |
| `make create-migration NAME=name`     | Create a new migration                        |
| `make migrate-help`                   | Display migration help                        |

## Graceful Shutdown

The application listens for:

```text
SIGINT
SIGTERM
```

When a shutdown signal is received:

1. The application context is cancelled.
2. Background processors are stopped.
3. The HTTP server is gracefully shut down.
4. Existing requests are given up to 10 seconds to complete.
5. Database connections are closed.

## CORS

CORS is currently configured to allow all origins during development.

For production, restrict allowed origins to the application's frontend domains.

## Development Workflow

A typical local development workflow:

```bash
# Install dependencies and development tools
make setup

# Create environment configuration
cp .env.sample .env

# Run database migrations
make migrate-up

# Generate SQLC code
make generate

# Seed sample data
make seed

# Start the API
go run .
```

The API should then be available at:

```text
http://localhost:8081
```

Health check:

```bash
curl http://localhost:8081/health
```

Use environment variables or your deployment platform's secret manager for production credentials.

## License

Add the project's license information here.
