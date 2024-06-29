package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"warehouse/api/warehousepb"
	"warehouse/internal/config"
	"warehouse/internal/db"
	intgrpc "warehouse/internal/grpc"
	articlesrepo "warehouse/internal/repositories/articles"
	productsrepo "warehouse/internal/repositories/products"
	"warehouse/internal/services/products"
)

func main() {
	fx.New(
		fx.Provide(NewApplicationContext),
		fx.Provide(config.NewConfig),
		fx.Provide(NewGRPCServer),
		fx.Provide(NewDatabase),
		fx.Provide(articlesrepo.NewRepository),
		fx.Provide(productsrepo.NewRepository),
		fx.Provide(NewWarehouseService),
		fx.Invoke(MigrateDatabase),
		fx.Invoke(func(server *grpc.Server, service *intgrpc.Service) {
			warehousepb.RegisterWarehouseServiceServer(server, service)
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

func NewGRPCServer(lc fx.Lifecycle, appCfg config.Config) (*grpc.Server, error) {
	var cfg intgrpc.Config
	err := appCfg.GetConfig("grpc", &cfg)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	reflection.Register(server)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", cfg.Address())
			if err != nil {
				return err
			}
			go func() {
				log.Printf("grpc server listening on %s", cfg.Address())
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

func NewDatabase(lc fx.Lifecycle, appCtx context.Context, appCfg config.Config) (*pgxpool.Pool, error) {
	cfg, err := db.ParseConfig(appCfg)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.New(appCtx, cfg.DSN())
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func MigrateDatabase(appCfg config.Config) error {
	cfg, err := db.ParseConfig(appCfg)
	if err != nil {
		return err
	}
	if !cfg.Migrations {
		return nil
	}

	m, err := migrate.New("file://db/migrations", cfg.DSN())
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func NewWarehouseService(aRepo articlesrepo.Repository, pRepo productsrepo.Repository) (*intgrpc.Service, error) {
	productsSrv := products.NewService(aRepo, pRepo)
	return intgrpc.NewService(productsSrv), nil
}
