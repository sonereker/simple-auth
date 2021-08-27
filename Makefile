GOOGLE_APIS_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/googleapis/googleapis 2> /dev/null)
GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2 2> /dev/null)
PG_VALIDATE_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/envoyproxy/protoc-gen-validate 2> /dev/null)
GO_INSTALLED := $(shell which go)

install-tools:
ifndef GO_INSTALLED
	$(error "go is not installed, please run 'brew install go'")
endif

	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/envoyproxy/protoc-gen-validate

generate: install-tools
	protoc --proto_path=proto \
   		--proto_path=$(GOOGLE_APIS_DIR) \
   		--proto_path=$(GRPC_GATEWAY_DIR) \
   		--proto_path=$(PG_VALIDATE_DIR)  \
		--go_out=internal/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/pb \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=internal/pb \
		--grpc-gateway_opt=paths=source_relative \
		--swagger_out=logtostderr=true:internal/pb \
		--validate_out="lang=go:internal/pb" \
		--validate_opt=paths=source_relative \
		proto/v1/*.proto
	cp internal/pb/v1/users.swagger.json www/swagger.json
test:
	go test -race ./...

lint:
	golint ./...
	go vet ./...

test-integration:
	go test -tags integration -v ./users

build:
	go build -o bin/grpc_server cmd/grpc-server/*.go
	go build -o bin/http_server cmd/http-server/*.go
