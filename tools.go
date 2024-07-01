//go:build tools
// +build tools

package warehouse

import (
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "go.uber.org/mock/mockgen"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
