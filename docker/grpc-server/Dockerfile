FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache bash

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ../.. .

RUN go build -o grpc_server cmd/grpc-server/main.go

WORKDIR /app

RUN cp /build/grpc_server .

COPY docker/wait-for-it.sh .
RUN chmod +x wait-for-it.sh

EXPOSE 8070

ENTRYPOINT ["./wait-for-it.sh","db:5432","--","/app/grpc_server"]
