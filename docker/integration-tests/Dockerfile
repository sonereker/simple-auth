FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY ../.. .

RUN go mod download

ENTRYPOINT ["go", "test", "-tags", "integration", "./internal/users", "-v"]
