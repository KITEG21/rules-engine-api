# Rules Engine API - Project Summary

## Short Description

A RESTful API that evaluates configurable business rules (defined in JSON/YAML) against real-time data payloads. Used for dynamic pricing, eligibility checks, validation workflows, and policy enforcement.

---

## Overview

The Rules Engine API allows developers to define business rules as data (not code), store them in a database, and evaluate them at runtime against any JSON payload. Rules support complex conditions with AND/OR logic, multiple operators, and nested field access.

**Use Cases:**
- Dynamic discount/pricing engines
- User eligibility verification
- Input validation pipelines
- Policy enforcement
- Workflow automation

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                   Client                                    │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │ HTTP/JSON
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              main.go                                        │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                    │
│  │   Config    │───▶│  Migration  │───▶│    DB       │                    │
│  │  (Viper)   │    │  (migrate)  │    │  (pgx/v5)   │                    │
│  └─────────────┘    └─────────────┘    └──────┬──────┘                    │
└─────────────────────────────────────────────────┼───────────────────────────┘
                                                  │
                                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         chi Router (HTTP Layer)                              │
│  GET/POST/PUT/DELETE /api/v1/rules                                         │
│  POST /api/v1/rules/evaluate                                                │
└────────────────────────────────┬────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        internal/api/handler.go                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                    │
│  │  DTOs       │───▶│ Validation  │───▶│  Business   │                    │
│  │ (request/   │    │(go-validator│    │   Logic     │                    │
│  │  response)  │    │  + Parser)  │    │             │                    │
│  └─────────────┘    └─────────────┘    └──────┬──────┘                    │
└─────────────────────────────────────────────────┼───────────────────────────┘
                                                  │
                          ┌───────────────────────┼───────────────────────┐
                          ▼                       ▼                       ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│   PostgreSQL     │  │  internal/rules  │  │    SQLC          │
│   (rules store) │  │  (evaluation    │  │ (type-safe DB   │
│                  │  │   engine)       │  │  queries)       │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

### Components

| Layer | Package | Responsibility |
|-------|---------|-----------------|
| Entry | `cmd/server` | App bootstrap, config, graceful shutdown |
| HTTP | `internal/api` | Routes, handlers, DTOs, validation |
| Engine | `internal/rules` | Parser (JSON/YAML), Evaluator, AST |
| Data | `internal/store` | SQLC-generated DB layer |
| Config | `internal/config` | Viper YAML loader |
| Migrate | `internal/migrate` | golang-migrate runner |

---

## Key Technical Decisions

| Decision | Rationale |
|----------|-----------|
| **Chi Router** | Lightweight, idiomatic Go router with middleware support |
| **pgx/v5** | High-performance PostgreSQL driver with connection pooling |
| **sqlc** | Type-safe SQL queries at compile time, no runtime reflection |
| **golang-migrate** | Version-controlled DB migrations, native Go |
| **Viper** | Unified config from YAML + env vars + flags |
| **go-validator** | Declarative validation tags on DTOs |
| **AST-based Engine** | Rules parsed into Abstract Syntax Tree for evaluation |
| **Definition as JSONB** | Store rule definitions as native JSON in PostgreSQL |

### Rule Engine Design

```
Rule Definition (JSON)     AST (Node Tree)          Evaluation
─────────────────────      ───────────────         ───────────
{                        Node{                     Input Data
  "logic": "AND",    ┌──► Logic: "AND"         ───────────
  "conditions": [    │    Conditions: [...]    │
    {"field":        │    }                      │
      "user.type",   │    ┌────────────────┐    │
      "eq",          │    │ Node{          │    │
      "premium"      │───►│ Field: "user   │───► true/false
    },               │    │ .type"         │    │
    {"field":        │    │ Operator: "eq" │    │
      "order.total", │    │ Value: "prem"  │    │
      "gte",         │    │ }              │    │
      100            │    └────────────────┘    │
  ]                                           │
}                                                │
```

---

## Challenges & Trade-offs

### Challenges

1. **JSON Serialization in PostgreSQL**
   - pgx returns JSONB as `[]byte`, which serializes as base64 in JSON
   - Solution: Custom `MarshalJSON()` in response layer

2. **Type Safety with Dynamic Rules**
   - Rule definitions are `any` type in DTOs
   - Solution: Validate with parser before saving to DB

3. **Operator Flexibility**
   - Need to support many operator types (comparison, string, collection)
   - Solution: Switch-based evaluator with extensible operator registry

4. **Nested Field Resolution**
   - Data can have arbitrary nested structures
   - Solution: Dot notation parser with recursive traversal

### Trade-offs

| Trade-off | Impact |
|-----------|--------|
| **Parser creates new instance per request** | Minor allocation overhead; acceptable for MVP |
| **Regex compiled on every evaluation** | Could cache for high-throughput; deferred optimization |
| **No rule caching in DB** | Rules are fetched per request; could add Redis later |
| **Single-node evaluator** | Not distributed; sufficient for single-region deployments |
| **No authentication** | Per `needs.txt` - can add JWT middleware later |

---

## API Surface

```
POST   /api/v1/rules          # Create rule
GET    /api/v1/rules          # List active rules
GET    /api/v1/rules/{id}     # Get rule by ID
PUT    /api/v1/rules/{id}     # Update rule
DELETE /api/v1/rules/{id}     # Delete rule
POST   /api/v1/rules/evaluate # Evaluate rules against data
```

---

## Example Evaluation

**Rule:**
```json
{
  "logic": "AND",
  "conditions": [
    { "field": "user.type", "operator": "eq", "value": "premium" },
    { "field": "order.total", "operator": "gte", "value": 100 }
  ]
}
```

**Request:**
```json
{
  "data": {
    "user": { "type": "premium" },
    "order": { "total": 150 }
  },
  "ruleIds": [1]
}
```

**Response:**
```json
{
  "results": [
    { "ruleId": "1", "matched": true }
  ]
}
```

---

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.25 |
| Router | chi/v5 |
| DB Driver | pgx/v5 |
| ORM/Queries | sqlc |
| Migrations | golang-migrate |
| Config | Viper |
| Validation | go-validator |
| YAML | go.yaml.in/yaml/v3 |
| Database | PostgreSQL |
