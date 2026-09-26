# Testing Standards

## Non-negotiables

- **Every service in `services/` must have unit tests.** A new service or a new
  service method without tests must not be merged.
- Unit tests never hit the network. Driven ports are replaced with in-package
  stubs/mocks defined in the `_test.go` file.
- Test files live next to the code: `services/shop_service_test.go`.

## Style

- Table-driven tests with named cases:

```go
tests := []struct {
    name    string
    id      domain.ShopID
    gateway stubShopGateway
    wantErr error
}{
    {name: "rejects non-positive id", id: 0, wantErr: domain.ErrInvalidInput},
    ...
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
```

- Assert error identity with `errors.Is`, never string comparison.
- Cover at minimum: happy path, every validation failure, and gateway error
  propagation.
- Use plain stdlib `testing`; no assertion frameworks unless already a dependency.
- Mocks/stubs are hand-written structs with function fields, e.g.:

```go
type stubShopGateway struct {
    listFn func(ctx context.Context) ([]domain.Shop, error)
}
func (s stubShopGateway) ListShops(ctx context.Context) ([]domain.Shop, error) {
    return s.listFn(ctx)
}
```

## Adapters

- Adapter tests (when added) use `httptest.Server`; they verify URL, method,
  auth header, DTO mapping, and status → domain error mapping.

## Commands

```sh
make check   # build + vet + fmt-check + lint + test (run before finishing any change)
```

Individual targets: `make build`, `make vet`, `make fmt`, `make fmt-check`,
`make lint` (installs pinned golangci-lint if needed), `make test`.

`make check` must pass before any change is considered done.
