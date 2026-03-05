FROM golang:1.25-alpine AS builder
WORKDIR /build

# Copy shared module
COPY shared/ /build/shared/

# Copy service
COPY write-api/go.mod write-api/go.sum ./write-api/
WORKDIR /build/write-api
RUN go mod download

WORKDIR /build
COPY write-api/ ./write-api/

WORKDIR /build/write-api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget
WORKDIR /app
COPY --from=builder /app /app/server
EXPOSE 8082
ENTRYPOINT ["/app/server"]
