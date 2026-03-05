FROM golang:1.25-alpine AS builder
WORKDIR /build

# Copy shared module
COPY shared/ /build/shared/

# Copy service
COPY read-api/go.mod read-api/go.sum ./read-api/
WORKDIR /build/read-api
RUN go mod download

WORKDIR /build
COPY read-api/ ./read-api/

WORKDIR /build/read-api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget
WORKDIR /app
COPY --from=builder /app /app/server
EXPOSE 8083
ENTRYPOINT ["/app/server"]
