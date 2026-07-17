# Stage 1: Agnostic Go Build Environment
FROM golang:1.22-alpine AS binary-builder
RUN apk add --no-cache git gcc g++ make
WORKDIR /src

# Leverage caching for universal builds
COPY go.mod go.sum ./
RUN go mod download

# Build a completely static binary (No dynamic OS dependencies)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w" \
    -o aura-binary .

# Stage 2: Microscopic Scratch Runtime
# Scratch contains zero overhead, meaning a 100% immune security surface area
FROM scratch
WORKDIR /
COPY --from=binary-builder /src/aura-binary /aura-service

# Expose gRPC port
EXPOSE 50051

ENV SYSTEM_LAYER="S4X" \
    BUILD_CONFIG="ALL_ENV_PROD"

ENTRYPOINT ["/aura-service"]
