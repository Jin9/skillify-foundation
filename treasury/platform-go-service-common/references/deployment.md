# Deployment

Dockerfile, GitLab CI pipeline, and runtime dependencies.

## Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata tini \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" > /etc/timezone
RUN adduser -D -u 10001 appuser
USER appuser
COPY --from=builder /app/server /server
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/server"]
```

Conventions: non-root (UID 10001), `tini` as init, Bangkok TZ, static binary.

## CI Pipeline (GitLab CI)

Stages: `validate` → `test` → `containerize` → `scan` → `deploy`

- **validate**: tag validation, Dockerfile lint
- **test**: `go test ./... -race -cover`, dependency scan, secret detection
- **containerize**: Docker build + push
- **scan**: SonarQube, Rapid7, Fortify, NexusIQ
- **deploy**: Update image tag in Kustomize repository

## Runtime Dependencies

| Package | Purpose |
|---|---|
| `github.com/gin-gonic/gin` | HTTP framework |
| `github.com/IBM/sarama` | Kafka client |
| `cloud.google.com/go/firestore` | Firestore client |
| `github.com/jackc/pgx/v5` | PostgreSQL driver |
| `github.com/redis/go-redis/v9` | Redis client |
| `github.com/caarlos0/env/v11` | Environment config parsing |
| `github.com/stretchr/testify` | Test assertions & mocks |
| `github.com/google/uuid` | UUID generation |
| `github.com/cenkalti/backoff/v4` | Exponential backoff retry |
