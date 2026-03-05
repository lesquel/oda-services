module github.com/lesquel/oda-gateway

go 1.25.7

require github.com/lesquel/oda-shared v0.0.0

require (
	github.com/go-chi/chi/v5 v5.2.5
	github.com/go-chi/cors v1.2.2
	github.com/golang-jwt/jwt/v5 v5.3.1
)

replace github.com/lesquel/oda-shared => ../shared
