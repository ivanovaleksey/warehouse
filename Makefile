.PHONY: generate-proto
generate-proto:
	protoc \
		-I . \
		--go_out=. \
		--go-grpc_out=. \
		api/warehouse.proto

.PHONY: generate
generate: generate-proto
	go generate ./...
