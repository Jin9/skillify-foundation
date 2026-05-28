# Dockerfile, Makefile, CI Pipeline

---

## 1. Dockerfile

Multi-platform build with private-module authentication via Docker secrets (never `ARG`/`ENV` for credentials):

```dockerfile
# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG GIT_COMMIT

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.mod go.sum ./

# Secret-mount keeps credentials out of image layers and history.
# Build with:
#   docker build \
#     --secret id=git_username,env=GIT_USERNAME \
#     --secret id=git_password,env=GIT_PASSWORD \
#     ...
RUN --mount=type=secret,id=git_username \
  --mount=type=secret,id=git_password \
  --mount=type=cache,target=/go/pkg/mod \
  git config --global url."https://$(cat /run/secrets/git_username):$(cat /run/secrets/git_password)@gitdev.devops.krungthai.com/".insteadOf "https://gitdev.devops.krungthai.com/" \
  && go mod download \
  && rm -f /root/.gitconfig

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -trimpath -a -installsuffix cgo \
  -ldflags "-s -w -extldflags -static -X main.commit=${GIT_COMMIT}" \
  -o /out/api .

FROM alpine:3.23

RUN addgroup -S app && adduser -S -G app -u 10001 app \
  && apk upgrade --no-cache \
  && apk add --no-cache tini ca-certificates tzdata \
  && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
  && echo "Asia/Bangkok" > /etc/timezone \
  && update-ca-certificates

WORKDIR /home/app
COPY --from=build --chown=app:app /out/api /home/app/api
USER app

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/home/app/api"]
```

Key conventions:

- `--mount=type=secret` for private-module authentication (never `ARG` for credentials).
- `--mount=type=cache` for Go module and build caches.
- Non-root user `app` (UID 10001).
- `tini` as init process.
- `Asia/Bangkok` timezone.
- Static binary (`CGO_ENABLED=0`, `-trimpath`, `-installsuffix cgo`).

---

## 2. Makefile targets

| Target | Purpose |
|---|---|
| `make all` | Default — runs `mod`, `lint`, `test` |
| `make help` | List all available targets |
| `make setup` | Install dev tools and git hooks |
| `make upgrade` | Upgrade devtools (golangci-lint, govulncheck, swagger, etc.) |
| `make mod` | `go fmt ./...` + `go mod tidy` |
| `make lint` | `golangci-lint run --fix` |
| `make vuln` | `govulncheck ./...` |
| `make test` | `go test -v -race -covermode=atomic -coverpkg=./... ./...` |
| `make test-integration` | `go test -v -race -tags=integration ./...` |
| `make coverage` | Test + open HTML coverage report |
| `make precommit` | `all` + `vuln` + `go mod verify` + `go vet` |
| `make ci` | `precommit` + `diff` (no uncommitted changes) |
| `make diff` | Fail if working tree is dirty |
| `make docker` | Build Docker image |
| `make run` | Run container with `.env` file |
| `make swagger` | Generate OpenAPI spec |
| `make openapi` | Serve Swagger UI on `:8910` |
| `make bump-version` | Bump semver (`make bump-version version=v0.0.1`) |
| `make doc` | Serve pkgsite on `:6060` |
| `make clean` | Remove build artifacts and caches |
| `make release-cloudrun` | Build, push, deploy to Cloud Run |

---

## 3. GitLab CI pipeline

Stages: `validate` → `test` → `containerize` → `scan` → `deploy`.

| Stage | What runs |
|---|---|
| `validate` | Tag validation (release branch), Dockerfile validation |
| `test` | Dependency check, secret detection, Go unit tests with coverage |
| `containerize` | Docker build + push |
| `scan` | SonarQube, Rapid7, Fortify, NexusIQ |
| `deploy` | Update image tag in Kustomize repository |
