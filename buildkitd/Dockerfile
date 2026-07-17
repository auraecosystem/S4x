# Build binary safely in a controlled Go container environment
FROM golang:1.22-alpine AS compiler
RUN apk add --no-cache git make gcc g++
WORKDIR /workspace

# Cache Go modules definitions 
COPY go.mod go.sum ./
RUN go mod download

# Build production binary without debugging overhead
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service-app .

# Deploy inside an absolute bare-minimum execution runtime
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates

# Import compiled app from compiler stage
COPY --from=compiler /workspace/service-app .

EXPOSE 50051
ENV CONTAINER_ENGINE="MOBY" LAYER="S4X"

CMD ["./service-app"]
