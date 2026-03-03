FROM golang:1.25.6-bookworm AS builder

WORKDIR /src

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -ldflags="-s -w" -o api ./cmd/api/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates=20230311+deb12u1 \
    tzdata=2025b-0+deb12u2 \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -u 10001 -r -s /usr/sbin/nologin appuser

RUN mkdir /data && chown appuser:appuser /data

WORKDIR /app

COPY --from=builder /src/api .

RUN chown appuser:appuser /app

VOLUME ["/data"]

USER appuser

EXPOSE 8080

ENTRYPOINT ["./api"]