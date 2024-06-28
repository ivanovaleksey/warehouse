package main

import (
	"context"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"warehouse/api/warehousepb"
)

func main() {
	fx.New(
		fx.Provide(NewApplicationContext),
		fx.Provide(NewGRPCServer),
		fx.Invoke(func(server *grpc.Server) {
			warehousepb.RegisterWarehouseServer(server, warehousepb.UnimplementedWarehouseServer{})
		}),
	).Run()
}

func NewApplicationContext(lc fx.Lifecycle) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})
	return ctx
}

func NewGRPCServer(lc fx.Lifecycle) (*grpc.Server, error) {
	server := grpc.NewServer()
	reflection.Register(server)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":8080")
			if err != nil {
				return err
			}
			go func() {
				log.Println("grpc server listening on :8080")
				err = server.Serve(lis)
				if err != nil {
					log.Printf("error serving grpc server: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.GracefulStop()
			return nil
		},
	})
	return server, nil
}
