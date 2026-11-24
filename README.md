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

`go.mod` records Go module versions; `go.sum` records their checksums.
