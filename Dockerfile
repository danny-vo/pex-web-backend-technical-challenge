FROM golang:1.15-alpine AS build_base

# Don't cache index locally
RUN apk add --no-cache git
RUN apk add --no-cache make
RUN apk add --no-cache gcc
RUN apk add --no-cache musl-dev

# Set working directory
WORKDIR /tmp/fibonacci-backend

# Populate module caches based on Go dependencies
COPY go.mod .
COPY go.sum .

# Get dependencies
RUN go mod download

# Get everything else
COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v ./...

# Build executable
RUN make build

# Second stage image only for executing
FROM alpine:3.12
RUN apk add --no-cache bash
RUN apk add ca-certificates
RUN apk add --no-cache curl

# Copy executable
COPY --from=build_base \
        /tmp/fibonacci-backend/out/fibonacci-backend/fibonacci_server \
        /app/fibonacci-backend

# Expose port
EXPOSE 8080

# Execute the binary generated
CMD ["/app/fibonacci-backend"]