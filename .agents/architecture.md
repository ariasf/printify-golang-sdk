# Hexagonal Architecture Rules

This SDK follows hexagonal (ports & adapters) architecture. These rules are mandatory.

## Layout

```
pkg/printify/
├── printify.go        # SDK facade: wiring only, no logic
├── domain/            # entities & value objects
├── ports/
│   ├── driven/        # outbound interfaces (what the SDK needs, e.g. API gateways)
│   └── driving/       # inbound interfaces (what the SDK offers to consumers)
├── services/          # ALL business logic lives here
└── adapters/
    └── rest/          # HTTP implementation of driven ports
```

## Dependency rule

Dependencies always point inward. Never the reverse.

```
adapters ──► ports/driven ──► domain
services ──► ports/driven ──► domain
printify.go ──► services, adapters (wiring only)
```

- `domain` imports nothing from this module.
- `ports` import only `domain` and stdlib (`context`).
- `services` import `domain` and `ports/driven`. NEVER `adapters`.
- `adapters` import `domain` and `ports/driven`. NEVER `services`.

## Layer responsibilities

### domain
- Entities, value objects, and domain error sentinels only.
- Invariant checks may live on entities (e.g. `Validate()`), but no orchestration,
  no I/O, no JSON tags for transport concerns — transport mapping belongs in adapters.

### ports
- **Interfaces only.** No implementations, no helper functions, no logic.
- Keep interfaces minimal: only the methods callers actually need (interface
  segregation). Prefer one small interface per resource over one giant interface.
- Every method takes `context.Context` as its first parameter.

### services
- All business logic: validation, orchestration, pagination handling, retries
  policy decisions, mapping between inputs and gateway calls.
- Each service implements a `ports/driving` interface and depends only on
  `ports/driven` interfaces (injected via constructor).
- Constructors: `NewXxxService(gw driven.XxxGateway) *XxxService`.

### adapters
- Translate driven port calls into HTTP requests against the Printify API.
- Own all transport concerns: URLs, JSON (de)serialization DTOs, auth headers,
  HTTP status → domain error mapping.
- Zero business logic. If you're tempted to add an `if` about business rules
  here, it belongs in a service.

## Adding a new resource (from the OpenAPI spec in api/openapi.yaml)

1. Add entities to `domain/`.
2. Add a driven port interface in `ports/driven/` (minimal methods).
3. Add a driving port interface in `ports/driving/`.
4. Implement the service in `services/` + unit tests (see testing.md).
5. Implement the REST adapter in `adapters/rest/` with DTOs private to the package.
6. Wire it in `printify.go`.
