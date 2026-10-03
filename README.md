# sly64

**sly64** is a lightweight and flexible DNS server. It supports name-based routing, DNS64, and multiple upstream resolvers with the ability to bind them to specific network interfaces.

## Build

Requires Go 1.27 or newer.

```sh
go build ./cmd/sly64
go test -race ./...
go vet ./...
```

## Dependencies

Go dependencies are updated as of 2026-10-03, including `github.com/buglloc/certifi` v0.9.7. Exact versions are recorded in `go.mod` and verified through `go.sum`.
