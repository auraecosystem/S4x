# syntax=docker/dockerfile:1

############################
# Stage 1: Build
############################
FROM golang:1.22-alpine AS binary-builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    gcc \
    g++ \
    make \
    clang \
    lld \
    cmake \
    ninja \
    pkgconf \
    nasm \
    yasm

WORKDIR /src

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build fully static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w" \
    -o /out/aura-service .

############################
# Stage 2: Runtime
############################
FROM scratch

WORKDIR /

# Copy binary
COPY --from=binary-builder /out/aura-service /aura-service

# Optional: copy CA certificates if your app makes HTTPS requests
COPY --from=binary-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose gRPC port
EXPOSE 50051

ENV SYSTEM_LAYER="S4X" \
    BUILD_CONFIG="ALL_ENV_PROD"

ENTRYPOINT ["/aura-service"]
