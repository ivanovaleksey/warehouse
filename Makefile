.PHONY: generate-proto
generate-proto:
	protoc \
		-I . \
		--go_out=. \
		--go-grpc_out=. \
		api/warehouse.proto

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: generate
generate: generate-proto gogen

.PHONY: test
test:
	go test -v -count 1 ./...

.PHONY: install-tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install go.uber.org/mock/mockgen
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: build-server
build-server:
	go build -o ./bin/server ./cmd/server

.PHONY: build-seed
build-seed:
	go build -o ./bin/seed ./cmd/seed
