package products

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"warehouse/internal/models"
)

var (
	ErrNotFound = errors.New("product not found")
)

type Repository interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProduct(ctx context.Context, id int32) (models.Product, error)
}

type impl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &impl{
		db: db,
	}
}

func (repo *impl) GetProducts(ctx context.Context) ([]models.Product, error) {
	const query = `
		SELECT id, name, price, articles
		FROM products
		ORDER BY id
	`

	var items []models.Product
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Product
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Articles)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		items = append(items, item)
	}
	return items, nil
}

func (repo *impl) GetProduct(ctx context.Context, id int32) (models.Product, error) {
	const query = `
		SELECT id, name, price, articles
		FROM products
		WHERE id = $1
	`
	var item models.Product
	err := repo.db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.Price, &item.Articles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, ErrNotFound
		}
		return models.Product{}, err
	}
	return item, nil
}
