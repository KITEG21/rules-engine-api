# Rules Engine API

A RESTful API for evaluating configurable rules in JSON/YAML format against real-time data. Useful for discount systems, validations, eligibility checks, and business logic enforcement.
<img width="1282" height="747" alt="image" src="https://github.com/user-attachments/assets/939788ad-9ac3-4473-96e1-dc50f9bba19f" />

## Features

- **CRUD Operations** - Create, read, update, delete rules
- **Rule Evaluation** - Evaluate rules against data payloads in real-time
- **JSON/YAML Support** - Define rules in either JSON or YAML format
- **Multiple Operators** - Support for comparisons, string operations, collections
- **Nested Logic** - AND/OR logic with nested conditions
- **Input Validation** - Request validation using go-validator
- **Graceful Shutdown** - Proper handling of server termination

## Tech Stack

| Component | Technology |
|-----------|------------|
| HTTP Router | Chi v5 |
| Database Driver | pgx/v5 |
| Query Generator | sqlc |
| Database Migrations | golang-migrate |
| Configuration | Viper |
| Validation | go-validator |
| YAML Parsing | go.yaml.in/yaml/v3 |

## Project Structure

```
rules_engine_api/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── dto/                    # Data Transfer Objects
│   │   │   ├── request.go          # Request DTOs
│   │   │   ├── response.go         # Response DTOs
│   │   │   └── validator.go        # Validation helper
│   │   ├── handler.go              # HTTP handlers
│   │   └── routes.go               # Route definitions
│   ├── config/
│   │   └── config.go               # Configuration loader
│   ├── migrate/
│   │   └── migrate.go              # Migration runner
│   ├── rules/
│   │   ├── node.go                 # AST Node definition
│   │   ├── parser.go               # JSON/YAML parser
│   │   ├── evaluator.go            # Rule evaluation engine
│   │   ├── operators.go            # Supported operators
│   │   ├── result.go               # Evaluation result
│   │   └── definition.go           # Definition type
│   └── store/                      # SQLC generated
│       ├── db.go                   # DB connection wrapper
│       ├── models.go               # Data models
│       └── queries.sql.go          # Query functions
├── migrations/
│   ├── 001_create_rules.up.sql
│   └── 001_create_rules.down.sql
├── config.yaml                     # Application config
├── sqlc.yaml                       # sqlc configuration
├── go.mod
└── README.md
```

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/rules` | List all active rules |
| POST | `/rules` | Create a new rule |
| GET | `/rules/{id}` | Get rule by ID |
| PUT | `/rules/{id}` | Update rule |
| DELETE | `/rules/{id}` | Delete rule |
| POST | `/rules/evaluate` | Evaluate rules against data |

---

## Request/Response Examples

### Create Rule

**Request**
```http
POST /api/v1/rules
Content-Type: application/json

{
  "name": "premium-discount",
  "description": "Apply 10% discount for premium users",
  "definition": {
    "logic": "AND",
    "conditions": [
      { "field": "user.type", "operator": "eq", "value": "premium" },
      { "field": "order.total", "operator": "gte", "value": 100 }
    ]
  }
}
```

**Response**
```json
{
  "id": 1,
  "name": "premium-discount",
  "description": "Apply 10% discount for premium users",
  "definition": {
    "logic": "AND",
    "conditions": [
      { "field": "user.type", "operator": "eq", "value": "premium" },
      { "field": "order.total", "operator": "gte", "value": 100 }
    ]
  },
  "is_active": true,
  "created_at": "2026-03-26T12:00:00.000000",
  "updated_at": "2026-03-26T12:00:00.000000"
}
```

### List Rules

**Request**
```http
GET /api/v1/rules
```

**Response**
```json
[
  {
    "id": 1,
    "name": "premium-discount",
    "description": "Apply 10% discount for premium users",
    "definition": { ... },
    "is_active": true,
    "created_at": "2026-03-26T12:00:00.000000",
    "updated_at": "2026-03-26T12:00:00.000000"
  }
]
```

### Evaluate Rules

**Request**
```http
POST /api/v1/rules/evaluate
Content-Type: application/json

{
  "data": {
    "user": { "type": "premium", "age": 25 },
    "order": { "total": 150, "items": ["item1", "item2"] }
  },
  "ruleIds": [1, 2]
}
```

**Response**
```json
{
  "results": [
    { "ruleId": "1", "matched": true },
    { "ruleId": "2", "matched": false, "error": "rule not found" }
  ]
}
```

---

## Rule Definition Syntax

### Simple Condition

```json
{
  "field": "user.age",
  "operator": "gte",
  "value": 18
}
```

### Logical Conditions (AND/OR)

```json
{
  "logic": "AND",
  "conditions": [
    { "field": "user.type", "operator": "eq", "value": "premium" },
    { "field": "order.total", "operator": "gte", "value": 100 }
  ]
}
```

### Nested Conditions

```json
{
  "logic": "OR",
  "conditions": [
    {
      "logic": "AND",
      "conditions": [
        { "field": "user.type", "operator": "eq", "value": "premium" },
        { "field": "order.total", "operator": "gte", "value": 500 }
      ]
    },
    { "field": "user.isVip", "operator": "eq", "value": true }
  ]
}
```

## Supported Operators

| Category | Operators | Description |
|----------|-----------|-------------|
| **Comparison** | `eq`, `equals` | Equal to |
| | `neq`, `not_equals` | Not equal to |
| | `gt` | Greater than |
| | `gte` | Greater than or equal |
| | `lt` | Less than |
| | `lte` | Less than or equal |
| **String** | `contains` | String contains substring |
| | `startswith` | String starts with prefix |
| | `endswith` | String ends with suffix |
| | `matches` | Regex match |
| **Collection** | `in` | Value in array |
| **Existence** | `exists` | Field exists |

## Field Path Resolution

Use dot notation to access nested fields in your data:

```json
{ "field": "user.profile.age" }
```

```json
{
  "user": {
    "profile": {
      "age": 25
    }
  }
}
```

## Setup

### Prerequisites

- Go 1.25+
- PostgreSQL 14+

### 1. Clone and Install Dependencies

```bash
cd rules_engine_api
go mod download
go mod tidy
```

### 2. Configure Database

Edit `config.yaml`:

```yaml
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "your_password"
  name: "rules_engine"

app:
  host: "0.0.0.0"
  port: 8080
```

### 3. Create Database

```bash
createdb rules_engine
```

### 4. Run Migrations (Automatic)

The app automatically runs migrations on startup. Alternatively:

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/rules_engine?sslmode=disable" up
```

### 5. Run the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Testing with cURL

```bash
# Create a rule
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{"name": "adult-check", "definition": {"field": "user.age", "operator": "gte", "value": 18}}'

# List rules
curl http://localhost:8080/api/v1/rules

# Evaluate rules
curl -X POST http://localhost:8080/api/v1/rules/evaluate \
  -H "Content-Type: application/json" \
  -d '{"data": {"user": {"age": 25}}, "ruleIds": [1]}'
```

## Error Responses

### Validation Error (400)

```json
{
  "error": "validation failed",
  "details": [
    { "field": "name", "message": "field is required" }
  ]
}
```

### Not Found (404)

```json
{
  "error": "rule not found"
}
```

### Server Error (500)

```json
{
  "error": "failed to create rule: ..."
}
```

## Configuration

The application supports multiple configuration sources (in order of precedence):

1. Command-line flags
2. Environment variables
3. config.yaml file
4. Default values

### Config File Format (config.yaml)

```yaml
database:
  host: string      # Database host (default: "localhost")
  port: int         # Database port (default: 5432)
  user: string      # Database user
  password: string  # Database password
  name: string      # Database name

app:
  host: string      # Server host (default: "0.0.0.0")
  port: int         # Server port (default: 8080)
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `DB_HOST` | Database host |
| `DB_PORT` | Database port |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_NAME` | Database name |
| `APP_HOST` | Server host |
| `APP_PORT` | Server port |
| `AI_API_KEY` | AI provider API key (for natural-language rule translation) |
| `AI_BASE_URL` | Optional AI base endpoint |
| `AI_MODEL` | AI model identifier (e.g., gemini-3-flash-preview) |

## AI Natural-Language Rule Translation

The API can accept rule definitions in natural language when AI is configured:

- `POST /api/v1/rules` with `definition` as string.
- The `internal/api/handler.go` converts plain text into AST using `internal/ai/client.go`.
- The AI client uses `config.yaml` values through `cmd/server/main.go`.
- If no AI key/model or translation fails, request returns 400 with a descriptive error (e.g. "failed to translate natural language definition: ...").

Example:

```json
{
  "name": "Adult in US",
  "description": "Check user is adult in US",
  "definition": "User is 18 or older and located in the United States"
}
```

## SHOWCASE / Project Summary

### Short Description

A RESTful API that evaluates configurable business rules (defined in JSON/YAML) against real-time data payloads. Used for dynamic pricing, eligibility checks, validation workflows, and policy enforcement.

### Supported Flows

- Dynamic discount/pricing engines
- User eligibility verification
- Input validation pipelines
- Policy enforcement
- Workflow automation

### Architecture (high-level)

- `cmd/server` bootstraps config, database migration, and API routes
- `internal/api` handles request parsing, validation, and rule CRUD/evaluate APIs
- `internal/rules` parses and evaluates rule AST definitions
- `internal/store` contains sqlc-generated persistence logic to PostgreSQL
- `internal/config` loads configuration with Viper (config file + env)
- `internal/ai` translates natural language rules to AST JSON via GenAI

### Rule Evaluation

Rules are stored as JSON definition (`definition` field) and parsed into AST nodes:

- `logic` (`AND` / `OR`)
- `conditions` array of nested nodes
- leaf nodes with `field`, `operator`, `value`

Evaluation returns `matched` boolean and optional `value`.

### Error responses

- 400: validation or translate error
- 404: rule not found
- 500: internal server error

## License

MIT

MIT
