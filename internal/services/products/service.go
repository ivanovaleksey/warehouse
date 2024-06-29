package products

import (
	"context"
	"fmt"

	"warehouse/internal/models"
	"warehouse/internal/repositories/articles"
	"warehouse/internal/repositories/products"
)

type Service interface {
	GetProductsWithStock(ctx context.Context) ([]models.ProductWithStock, error)
	RemoveProduct(ctx context.Context, id, quantity int32) error
}

type impl struct {
	articlesRepo articles.Repository
	productsRepo products.Repository
}

func NewService(aRepo articles.Repository, pRepo products.Repository) Service {
	return &impl{
		articlesRepo: aRepo,
		productsRepo: pRepo,
	}
}

// GetProductsWithStock list the products and calculates stock quantity
func (srv *impl) GetProductsWithStock(ctx context.Context) ([]models.ProductWithStock, error) {
	prods, err := srv.productsRepo.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	arts, err := srv.articlesRepo.GetArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %w", err)
	}

	inventory := make(map[int32]int32, len(arts))
	for _, article := range arts {
		inventory[article.ID] = article.Stock
	}

	prodsWithStock := make([]models.ProductWithStock, 0, len(prods))
	for _, prod := range prods {
		minStock := int32(0)
		for i, art := range prod.Articles {
			artStock, ok := inventory[art.ID]
			if !ok {
				minStock = 0
				break
			}

			stock := artStock / art.Quantity
			if i == 0 {
				minStock = stock
			}
			if stock < minStock {
				minStock = stock
			}
		}
		prodsWithStock = append(prodsWithStock, models.ProductWithStock{
			Product: prod,
			Stock:   minStock,
		})
	}
	return prodsWithStock, nil
}

// RemoveProduct checks available quantity on stock and removes the product from stock
// TODO: check available quantity
func (srv *impl) RemoveProduct(ctx context.Context, id, quantity int32) error {
	product, err := srv.productsRepo.GetProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}
	if len(product.Articles) == 0 {
		return nil
	}

	arts := make([]models.ProductArticle, 0, len(product.Articles))
	for _, a := range product.Articles {
		arts = append(arts, models.ProductArticle{
			ID:       a.ID,
			Quantity: a.Quantity * quantity,
		})
	}
	return srv.articlesRepo.RemoveArticles(ctx, arts)
}
