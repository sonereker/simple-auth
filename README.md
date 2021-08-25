# Simple-Auth

[![Go Report Card](https://goreportcard.com/badge/github.com/sonereker/simple-auth)](https://goreportcard.com/report/github.com/sonereker/simple-auth)

An example gRPC/REST/JWT user authentication service.

## Quick Run

```
docker-compose build
docker-compose up
```

## Development

Install required tools;

```
make install-tools
```

Generate gRPC and REST bindings;

```
make generate
```

Environment variables required;

```
DB_HOST=localhost;DB_NAME=simple_auth;DB_USERNAME=local;DB_PASSWORD=local;DB_PORT=5432;DB_SSL_MODE=disable
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
go test -tags integration -v ./..
```
