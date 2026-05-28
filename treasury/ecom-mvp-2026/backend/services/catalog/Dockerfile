# syntax=docker/dockerfile:1.7
# Multi-stage build: minimal runtime image, expose 8080, healthcheck
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG GIT_COMMIT

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy go workspace files first for layer caching
COPY go.mod go.sum ./

# Download dependencies
# The replace directive in go.mod points to ../../common (monorepo local path).
# In CI, copy common alongside this service before building, or use a Git tag replace.
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -a -installsuffix cgo \
    -ldflags "-s -w -extldflags -static -X main.commit=${GIT_COMMIT}" \
    -o /out/catalog .

# ---------- Runtime stage ----------
FROM alpine:3.23

RUN addgroup -S app && adduser -S -G app -u 10001 app \
    && apk upgrade --no-cache \
    && apk add --no-cache tini ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" > /etc/timezone \
    && update-ca-certificates

WORKDIR /home/app

COPY --from=build --chown=app:app /out/catalog /home/app/catalog

USER app

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/liveness || exit 1

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/home/app/catalog"]
