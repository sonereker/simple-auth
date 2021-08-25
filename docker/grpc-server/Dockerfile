FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ../.. .

RUN go build -o grpc_server cmd/grpc-server/main.go

WORKDIR /app

RUN cp /build/grpc_server .

EXPOSE 8070

CMD ["/app/grpc_server"]
