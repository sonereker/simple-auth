FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o auth cmd/server/main.go

WORKDIR /app

RUN cp /build/auth .

EXPOSE 8080

CMD ["/app/auth"]
