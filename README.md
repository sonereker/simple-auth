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

## Using Client

Client will register a new user, which will also login user.
```
make demo
```
