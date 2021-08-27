FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY /go.mod .
COPY /go.sum .
RUN go mod download

COPY ../.. .

RUN go build -o http_server cmd/http-server/main.go

WORKDIR /app

RUN cp /build/http_server .

COPY /www www/
EXPOSE 8080

CMD ["/app/http_server"]
