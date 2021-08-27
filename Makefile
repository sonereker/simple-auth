GRPC_GATEWAY_DIR := $(shell go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway 2> /dev/null)
GO_INSTALLED := $(shell which go)
PROTOC_INSTALLED := $(shell which protoc)
BINDATA_INSTALLED := $(shell which go-bindata 2> /dev/null)
PGGG_INSTALLED := $(shell which protoc-gen-grpc-gateway 2> /dev/null)
PGG_INSTALLED := $(shell which protoc-gen-go 2> /dev/null)

install-tools:
ifndef GO_INSTALLED
	$(error "go is not installed, please run 'brew install go'")
endif
ifndef PROTOC_INSTALLED
	$(error "protoc is not installed, please run 'brew install protobuf'")
endif
ifndef BINDATA_INSTALLED
	@go get -u github.com/kevinburke/go-bindata/go-bindata@master
endif
ifndef PGGG_INSTALLED
	@go get -u github.com/grpc-ecosystem/grpc-gateway/...
endif
ifndef PGG_INSTALLED
	@go get -u github.com/golang/protobuf/protoc-gen-go
endif

generate: install-tools
	protoc --proto_path=proto \
   		--proto_path=$(GRPC_GATEWAY_DIR)/third_party/googleapis \
		--go_out=internal/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/pb \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=internal/pb \
		--grpc-gateway_opt=paths=source_relative \
		--swagger_out=logtostderr=true:internal/pb \
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
