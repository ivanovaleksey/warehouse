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
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"warehouse/api/warehousepb"
	"warehouse/internal/config"
	"warehouse/internal/db"
	intgrpc "warehouse/internal/grpc"
)

func main() {
	fx.New(
		fx.Provide(NewApplicationContext),
		fx.Provide(NewConfig),
		fx.Provide(NewGRPCServer),
		fx.Provide(NewDatabase),
		fx.Invoke(MigrateDatabase),
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

func NewDatabase(lc fx.Lifecycle, appCfg config.Config) (*pgxpool.Pool, error) {
	var cfg db.Config
	err := appCfg.GetConfig("database", &cfg)
	if err != nil {
		return nil, err
	}
	pgxCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
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
	const source = "file://db/migrations"

	var cfg db.Config
	err := appCfg.GetConfig("database", &cfg)
	if err != nil {
		return err
	}
	if !cfg.Migrations {
		return nil
	}

	m, err := migrate.New(source, cfg.DSN())
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func NewConfig() config.Config {
	k := koanf.New(".")
	f := file.Provider("config.local.yaml")
	if err := k.Load(f, yaml.Parser()); err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
	return config.NewLoader(k)
}
