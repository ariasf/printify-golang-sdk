# Go Style Guide

- Target the Go version declared in go.mod. Run `gofmt` (or `go fmt ./...`) on
  every change.
- `context.Context` is always the first parameter of any method that can block
  or do I/O. Never store contexts in structs.
- Errors:
  - Sentinel errors live in `domain/errors.go` (`domain.ErrNotFound`, ...).
  - Wrap with context: `fmt.Errorf("listing shops: %w", err)`.
  - Callers check with `errors.Is` / `errors.As`.
- Accept interfaces, return concrete types. Ports are defined where they are
  consumed conceptually (`ports/`), not next to implementations.
- Exported identifiers need doc comments starting with the identifier name.
- No global state, no `init()` side effects, no panics in library code.
- Typed IDs over bare primitives: `type ShopID int64`, not `int64`.
- DTO structs with JSON tags are private to `adapters/rest`; domain entities
  carry no transport tags.
- Keep the public surface small: consumers interact via `printify.New(...)` and
  the `ports/driving` interfaces only.
- Dependencies: stdlib first. Adding a third-party dependency requires a strong
  justification.
