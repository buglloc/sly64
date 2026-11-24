ARG BASE=gcr.io/distroless/static-debian13:nonroot

FROM docker.io/library/golang:1.27.4 AS build

SHELL [ "/bin/sh", "-ec" ]
WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/sly64 ./cmd/sly64

FROM ${BASE}

COPY --from=build /go/bin/sly64 /sly64
USER nonroot:nonroot

WORKDIR /
EXPOSE 5553 53/udp
EXPOSE 5553 53/tcp

ENTRYPOINT ["/sly64"]
