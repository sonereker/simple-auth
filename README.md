# Simple Auth

An example user registration and authentication service featuring;
- gRPC Server
- HTTP Server with REST endpoints generated with `grpc-gateway`
- JWT for authentication
- Swagger REST API documentation
- Containerized services/integration tests

## Quick Run

```
docker-compose build
docker-compose up
```

## REST API Documentation

Swagger UI is available (see `Quick Run`) at http://localhost:8080/swagger-ui/


## Development

Install required tools;

```
make install-tools
```

Generate gRPC and REST bindings;

```
make generate
```

Environment variables required for database connection;
```
DB_HOST=localhost;DB_NAME=simple_auth;DB_USERNAME=local;DB_PASSWORD=local;DB_PORT=5432;DB_SSL_MODE=disable
```

For integration tests `GRPC_SERVER_ADDR` variable is also required. See `docker-compose.yml` for details.

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

or easier way is just running `integration-test` with docker-compose;

```
docker-compose build
docker-compose up
```
