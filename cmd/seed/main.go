package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"go.uber.org/fx"

	"warehouse/internal/config"
	"warehouse/internal/db"
)

func main() {
	fx.New(
		fx.Provide(NewApplicationContext),
		fx.Provide(config.NewConfig),
		fx.Provide(NewDatabase),
		fx.Invoke(func(sd fx.Shutdowner, appCtx context.Context, db *pgx.Conn) error {
			err := SeedProducts(appCtx, db)
			if err != nil {
				return err
			}
			err = SeedArticles(appCtx, db)
			if err != nil {
				return err
			}

			return sd.Shutdown()
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

func NewDatabase(lc fx.Lifecycle, appCtx context.Context, appCfg config.Config) (*pgx.Conn, error) {
	cfg, err := db.ParseConfig(appCfg)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(appCtx, cfg.DSN())
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close(ctx)
		},
	})

	return conn, nil
}

func SeedProducts(ctx context.Context, db *pgx.Conn) error {
	f, err := os.Open("cmd/seed/data/products.json")
	if err != nil {
		return err
	}
	defer f.Close()

	var content struct {
		Products []struct {
			Name            string `json:"name"`
			ContainArticles []struct {
				ArtId    string `json:"art_id"`
				AmountOf string `json:"amount_of"`
			} `json:"contain_articles"`
			Price int `json:"price"`
		} `json:"products"`
	}
	err = json.NewDecoder(f).Decode(&content)
	if err != nil {
		return err
	}

	table := "products"
	columns := []string{"name", "price", "articles"}
	var rows [][]any
	for _, item := range content.Products {
		var articles []any
		for _, a := range item.ContainArticles {
			id, err := strconv.Atoi(a.ArtId)
			if err != nil {
				return err
			}
			quantity, err := strconv.Atoi(a.AmountOf)
			if err != nil {
				return err
			}
			articles = append(articles, struct {
				ID       int32
				Quantity int32
			}{
				ID:       int32(id),
				Quantity: int32(quantity),
			})
		}
		rows = append(rows, []any{item.Name, item.Price, articles})
	}
	_, err = db.CopyFrom(ctx, pgx.Identifier{table}, columns, pgx.CopyFromRows(rows))
	return err
}

func SeedArticles(ctx context.Context, db *pgx.Conn) error {
	f, err := os.Open("cmd/seed/data/inventory.json")
	if err != nil {
		return err
	}
	defer f.Close()

	var content struct {
		Articles []struct {
			ArtId string `json:"art_id"`
			Name  string `json:"name"`
			Stock string `json:"stock"`
		} `json:"inventory"`
	}
	err = json.NewDecoder(f).Decode(&content)
	if err != nil {
		return err
	}

	table := "articles"
	columns := []string{"id", "name", "stock"}
	var rows [][]any
	for _, item := range content.Articles {
		rows = append(rows, []any{item.ArtId, item.Name, item.Stock})
	}
	_, err = db.CopyFrom(ctx, pgx.Identifier{table}, columns, pgx.CopyFromRows(rows))
	return err
}
