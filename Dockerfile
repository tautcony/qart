# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Setup for proxy
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.io,direct

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o qart .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/qart .

# Copy necessary files
COPY --from=builder /build/conf ./conf
COPY --from=builder /build/views ./views
COPY --from=builder /build/static ./static

# Create storage directories
RUN mkdir -p storage/flag storage/qrsave

EXPOSE 8080

ENTRYPOINT ["./qart"]
CMD ["--prod", "--port", "8080"]
