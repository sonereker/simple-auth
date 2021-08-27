# Simple Auth

[![Go](https://github.com/sonereker/simple-auth/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/sonereker/simple-auth/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/sonereker/simple-auth?status.png)](http://godoc.org/github.com/sonereker/simple-auth)
[![Go Report Card](https://goreportcard.com/badge/github.com/sonereker/simple-auth?2)](https://goreportcard.com/report/github.com/sonereker/simple-auth)

An example user registration and authentication service featuring;

- gRPC Server
- HTTP Server with RESTful endpoints
- JWT-based authentication
- Swagger REST API documentation
- Containerized integration tests

Service uses following Go packages;

- Database ORM: [gorm](https://github.com/go-gorm/gorm)
- JWT: [golang-jwt](https://github.com/golang-jwt/jwt)
- RESTful Endpoints: [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

## Quick Run

```
docker-compose build
docker-compose up
```

will build and start these containers;

- `grpc_server`
- `http_server`
- `db` (postgres:13.1-alpine)
- `integration-tests`

## RESTful API

After `Quick Run`, you have a RESTful API server running at http://localhost:8080. It provides the following endpoints:

- `POST /v1/users`: create a new user
- `POST /v1/users/login`: login with given credentials
- `GET /v1/users/current`: get current user

with cURL:

1. Create a new user
```
curl --location --request POST 'localhost:8080/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "dummy@email.com",
    "password": "hello123"
}'
```

2. Login with given credentials
```
curl --location --request POST 'localhost:8080/v1/users/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "dummy@email.com",
    "password": "hello123"
}'
```

3. Get current user
```
curl --location --request GET 'localhost:8080/v1/users/current' \
--header 'Authorization: <TOKEN_FROM_STEP2>'
```

### API Documentation

Swagger UI is available (see `Quick Run`) at http://localhost:8080/swagger-ui/

## Development

Install required tools:

```
make install-tools
```

Generate gRPC and REST bindings:

```
make generate
```

Environment variables required when running services locally:

```
DB_HOST=localhost;
DB_NAME=simple_auth;
DB_USERNAME=local;
DB_PASSWORD=local;
DB_PORT=5432;
DB_SSL_MODE=disable

GRPC_SERVER_ADDR=localhost:8070
HTTP_SERVER_ADDR=localhost:8080
```

## Tests

### Unit Tests

```
go test -v ./..
```

### Integration Tests

There's the `users/service_integration_test.go` test covering the basic functionality. After running gRPC server and
database server (see Quick Run) just run test with `tags` flag.

```
make test-integation
```

or easier way is just running `integration-test` with docker-compose:

```
docker-compose build
docker-compose up
```
