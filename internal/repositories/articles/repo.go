package articles

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"warehouse/internal/models"
)

var (
	ErrNotFound = errors.New("article not found")
)

type Repository interface {
	GetArticles(ctx context.Context) ([]models.Article, error)
	GetArticle(ctx context.Context, id int32) (models.Article, error)
	RemoveArticles(ctx context.Context, items []models.ProductArticle) error
}

type impl struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &impl{
		db: db,
	}
}

func (repo *impl) GetArticles(ctx context.Context) ([]models.Article, error) {
	const query = `
		SELECT id, name, stock
		FROM articles
		ORDER BY id
	`

	var items []models.Article
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Article
		err := rows.Scan(&item.ID, &item.Name, &item.Stock)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		items = append(items, item)
	}
	return items, nil
}

func (repo *impl) GetArticle(ctx context.Context, id int32) (models.Article, error) {
	const query = `
		SELECT id, name, stock
		FROM articles
		WHERE id = $1
	`
	var item models.Article
	err := repo.db.QueryRow(ctx, query, id).Scan(&item.ID, &item.Name, &item.Stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Article{}, ErrNotFound
		}
		return models.Article{}, err
	}
	return item, nil
}

func (repo *impl) RemoveArticles(ctx context.Context, items []models.ProductArticle) error {
	const query = `
		WITH to_remove (id, quantity) AS (
			SELECT *
			FROM unnest($1::int[], $2::int[])
		)
		UPDATE articles
		SET stock = stock - to_remove.quantity
		FROM to_remove
		WHERE articles.id = to_remove.id
	`

	ids := make([]int32, len(items))
	stocks := make([]int32, len(items))
	for i, item := range items {
		ids[i] = item.ID
		stocks[i] = item.Quantity
	}

	_, err := repo.db.Exec(ctx, query, ids, stocks)
	return err
}
