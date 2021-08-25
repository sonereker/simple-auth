# Simple-Auth

An example REST/gRPC authentication service.

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

## Using Client

Client will register a new user, which will also login user.
```
make demo
```
