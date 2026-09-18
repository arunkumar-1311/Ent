# Shipment Service

A simple REST API for managing shipments, built with [Echo](https://echo.labstack.com/) (HTTP framework) and [Ent](https://entgo.io/) (entity framework/ORM) on top of PostgreSQL.

## What This Project Does

The service exposes CRUD endpoints for a `Shipment` resource — tracking number, sender/receiver, origin/destination, weight, and status. It follows a layered architecture:

```
cmd/server      → application entry point (wires everything together, starts the HTTP server)
handler/        → HTTP layer (Echo handlers, request/response JSON)
service/        → business logic layer
repository/     → data access layer (talks to the database via Ent)
model/          → request/response DTOs
config/         → database connection setup
ent/            → Ent-generated ORM code + schema definitions
```

Request flow: `handler` → `service` → `repository` → `ent` (database).

## What is Ent?

[Ent](https://entgo.io/) is an entity framework for Go, developed by Facebook/Meta. Instead of hand-writing SQL or struct-based ORM mappings, you define your data model once as a Go schema, and Ent code-generates a fully-typed client, query builder, and migration logic for you.

In this project:

- The schema is defined in [ent/schema/shipment.go](ent/schema/shipment.go).
- Everything else under [ent/](ent/) (client, mutation builders, query builders, generated types, predicates, etc.) is **auto-generated** — you should never edit those files by hand.
- The `Shipment` entity currently has these fields:

  | Field             | Type    | Notes                        |
  |-------------------|---------|------------------------------|
  | `tracking_number` | string  | unique                       |
  | `sender_name`     | string  |                              |
  | `receiver_name`   | string  |                              |
  | `origin`          | string  |                              |
  | `destination`     | string  |                              |
  | `status`          | string  | defaults to `"pending"`      |
  | `weight`          | float   |                              |

- On every application startup ([cmd/server/main.go](cmd/server/main.go)), Ent's `client.Schema.Create(ctx)` runs an auto-migration, creating/updating the database schema (table `shipments` inside the `logistics` Postgres schema) to match the entity definition.
- Whenever you change [ent/schema/shipment.go](ent/schema/shipment.go) (add a field, add an edge/relation, etc.), you must regenerate the Ent code (see [Regenerating Ent Code](#regenerating-ent-code) below) before the change takes effect in the rest of the app.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26 or later
- [PostgreSQL](https://www.postgresql.org/download/) running locally (or reachable) with a database available
- (Optional) [`entc`](https://entgo.io/docs/code-gen) code generator, only needed if you plan to change the schema

## Setup — Step by Step

### 1. Clone / open the project

Make sure you're in the project root (the folder containing `go.mod`).

### 2. Install Go dependencies

```powershell
go mod download
```

This pulls in Echo, Ent, the Postgres driver (`lib/pq`), and their transitive dependencies as listed in [go.mod](go.mod).

### 3. Create the PostgreSQL database and schema

The app expects:
- A database named `rte`
- A Postgres **schema** (namespace) named `logistics` inside that database

```sql
-- Connect to Postgres as a superuser, then run:
CREATE DATABASE rte;
\c rte
CREATE SCHEMA IF NOT EXISTS logistics;
```

### 4. Configure the database connection

The connection string is currently hardcoded in [config/database.go](config/database.go):

```go
dsn := "host=localhost port=5432 user=postgres password=password dbname=rte sslmode=disable search_path=logistics"
```

Update the `host`, `port`, `user`, `password`, and `dbname` values there to match your local Postgres setup before running the app.

> Note: for production use, this DSN should be moved to environment variables instead of being hardcoded — it currently contains plaintext credentials in source code.

### 5. Run the server

```powershell
go run ./cmd/server
```

On startup the app will:
1. Connect to Postgres using the DSN above.
2. Auto-create/migrate the `shipments` table under the `logistics` schema.
3. Start the HTTP server on **port 8080**.

You should see the Echo server start without errors. If the database/schema doesn't exist or credentials are wrong, the app will log a fatal error and exit.

### 6. Verify it's running

```powershell
curl http://localhost:8080/shipments
```

This should return `[]` (empty list) on a fresh database.

## API Endpoints

| Method | Path              | Description              |
|--------|-------------------|--------------------------|
| POST   | `/shipments`      | Create a new shipment    |
| GET    | `/shipments`      | List all shipments       |
| GET    | `/shipments/:id`  | Get a shipment by ID     |
| PUT    | `/shipments/:id`  | Update a shipment by ID  |
| DELETE | `/shipments/:id`  | Delete a shipment by ID  |

### Example: Create a shipment

```powershell
curl -X POST http://localhost:8080/shipments `
  -H "Content-Type: application/json" `
  -d '{"tracking_number":"TRK123","sender_name":"Alice","receiver_name":"Bob","origin":"NYC","destination":"LA","weight":12.5}'
```

### Example: Update a shipment

```powershell
curl -X PUT http://localhost:8080/shipments/1 `
  -H "Content-Type: application/json" `
  -d '{"sender_name":"Alice","receiver_name":"Bob","origin":"NYC","destination":"LA","status":"shipped","weight":12.5}'
```

### Example: Delete a shipment

```powershell
curl -X DELETE http://localhost:8080/shipments/1
```

## Regenerating Ent Code

If you edit [ent/schema/shipment.go](ent/schema/shipment.go) (e.g., add a new field or entity), regenerate the Ent client code:

```powershell
go generate ./ent
```

This updates all the generated files under [ent/](ent/) to reflect the new schema. Restart the server afterward — the next `client.Schema.Create(ctx)` call will apply the corresponding database migration automatically.

## Ent Setup From Scratch

How the Ent layer was bootstrapped in this project (useful if setting up a new entity from zero):

### 1. Install Ent and scaffold an entity

```bash
go get entgo.io/ent/cmd/ent
go run entgo.io/ent/cmd/ent new Shipment
```

This generates `ent/generate.go` and `ent/schema/shipment.go`.

### 2. Define the schema fields

In [ent/schema/shipment.go](ent/schema/shipment.go):

```go
func (Shipment) Fields() []ent.Field {
	return []ent.Field{
		field.String("tracking_number").Unique(),
		field.String("sender_name"),
		field.String("receiver_name"),
		field.String("origin"),
		field.String("destination"),
		field.String("status").Default("pending"),
		field.Float("weight"),
	}
}
```

### 3. Generate the client code

```bash
go generate ./ent
```

or directly:

```bash
go run entgo.io/ent/cmd/ent generate ./ent/schema
```

This produces the typed `ent.Client`, query/mutation builders, and predicates used throughout [repository/shipment_repository.go](repository/shipment_repository.go).

### 4. Open the connection and auto-migrate

Done once in [config/database.go](config/database.go) via `ent.Open("postgres", dsn)`, and applied in [cmd/server/main.go](cmd/server/main.go) via `client.Schema.Create(ctx)`, which creates/updates the `shipments` table to match the schema.
