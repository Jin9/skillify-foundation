# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy the common library (replace directive points here at build time)
COPY ../../common /common

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -trimpath -a -installsuffix cgo \
  -ldflags "-s -w" \
  -o /out/identity .

FROM alpine:3.20

RUN addgroup -S app && adduser -S -G app -u 10001 app \
  && apk upgrade --no-cache \
  && apk add --no-cache tini ca-certificates tzdata \
  && cp /usr/share/zoneinfo/UTC /etc/localtime \
  && echo "UTC" > /etc/timezone \
  && update-ca-certificates

WORKDIR /home/app

COPY --from=build --chown=app:app /out/identity /home/app/identity

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/liveness || exit 1

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/home/app/identity"]
