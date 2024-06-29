package products

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"warehouse/internal/models"
	"warehouse/internal/repositories/articles/mock"
	"warehouse/internal/repositories/products/mock"
	"warehouse/internal/testhelpers"
)

func TestImpl_GetProductsWithStock(t *testing.T) {
	t.Run("basic scenario", func(t *testing.T) {
		fx := newFixture(t)

		products := []models.Product{
			{
				ID:    testhelpers.RandomInt32(),
				Name:  testhelpers.RandomString(),
				Price: float32(testhelpers.RandomInt32()),
				Articles: []models.ProductArticle{
					{
						ID:       1,
						Quantity: 4,
					},
					{
						ID:       2,
						Quantity: 8,
					},
					{
						ID:       3,
						Quantity: 1,
					},
				},
			},
			{
				ID:    testhelpers.RandomInt32(),
				Name:  testhelpers.RandomString(),
				Price: float32(testhelpers.RandomInt32()),
				Articles: []models.ProductArticle{
					{
						ID:       1,
						Quantity: 4,
					},
					{
						ID:       2,
						Quantity: 8,
					},
					{
						ID:       4,
						Quantity: 1,
					},
				},
			},
		}
		fx.productsRepo.EXPECT().GetProducts(fx.ctx).Return(products, nil)

		articles := []models.Article{
			{
				ID:    1,
				Stock: 12,
			},
			{
				ID:    2,
				Stock: 17,
			},
			{
				ID:    3,
				Stock: 2,
			},
			{
				ID:    4,
				Stock: 1,
			},
		}
		fx.articlesRepo.EXPECT().GetArticles(fx.ctx).Return(articles, nil)

		items, err := fx.GetProductsWithStock(fx.ctx)

		require.NoError(t, err)
		expected := []models.ProductWithStock{
			{
				Product: products[0],
				Stock:   2,
			},
			{
				Product: products[1],
				Stock:   1,
			},
		}
		assert.Equal(t, expected, items)
	})

	t.Run("should not fail if article is not in stock", func(t *testing.T) {
		fx := newFixture(t)

		products := []models.Product{
			{
				ID:    testhelpers.RandomInt32(),
				Name:  testhelpers.RandomString(),
				Price: float32(testhelpers.RandomInt32()),
				Articles: []models.ProductArticle{
					{
						ID:       2,
						Quantity: 2,
					},
					{
						ID:       1,
						Quantity: 4,
					},
				},
			},
		}
		fx.productsRepo.EXPECT().GetProducts(fx.ctx).Return(products, nil)

		articles := []models.Article{
			{
				ID:    1,
				Stock: 12,
			},
			{
				ID:    2,
				Stock: 1,
			},
		}
		fx.articlesRepo.EXPECT().GetArticles(fx.ctx).Return(articles, nil)

		items, err := fx.GetProductsWithStock(fx.ctx)

		require.NoError(t, err)
		expected := []models.ProductWithStock{
			{
				Product: products[0],
				Stock:   0,
			},
		}
		assert.Equal(t, expected, items)
	})

	t.Run("should not fail if article is unknown", func(t *testing.T) {
		fx := newFixture(t)

		products := []models.Product{
			{
				ID:    testhelpers.RandomInt32(),
				Name:  testhelpers.RandomString(),
				Price: float32(testhelpers.RandomInt32()),
				Articles: []models.ProductArticle{
					{
						ID:       1,
						Quantity: 4,
					},
					{
						ID:       3,
						Quantity: 1,
					},
				},
			},
		}
		fx.productsRepo.EXPECT().GetProducts(fx.ctx).Return(products, nil)

		articles := []models.Article{
			{
				ID:    1,
				Stock: 12,
			},
		}
		fx.articlesRepo.EXPECT().GetArticles(fx.ctx).Return(articles, nil)

		items, err := fx.GetProductsWithStock(fx.ctx)

		require.NoError(t, err)
		expected := []models.ProductWithStock{
			{
				Product: products[0],
				Stock:   0,
			},
		}
		assert.Equal(t, expected, items)
	})
}

func TestImpl_RemoveProduct(t *testing.T) {
	productID := testhelpers.RandomInt32()
	product := models.Product{
		Articles: []models.ProductArticle{
			{
				ID:       testhelpers.RandomInt32(),
				Quantity: int32(testhelpers.RandomIntRange(1, 100)),
			},
			{
				ID:       testhelpers.RandomInt32(),
				Quantity: int32(testhelpers.RandomIntRange(1, 100)),
			},
		},
	}

	t.Run("should not remove articles if product does not have any", func(t *testing.T) {
		fx := newFixture(t)

		product := models.Product{}
		fx.productsRepo.EXPECT().GetProduct(fx.ctx, productID).Return(product, nil)

		err := fx.RemoveProduct(fx.ctx, productID, 1)

		require.NoError(t, err)
	})

	t.Run("should remove 1 product", func(t *testing.T) {
		fx := newFixture(t)

		fx.productsRepo.EXPECT().GetProduct(fx.ctx, productID).Return(product, nil)
		articles := []models.ProductArticle{
			{
				ID:       product.Articles[0].ID,
				Quantity: product.Articles[0].Quantity,
			},
			{
				ID:       product.Articles[1].ID,
				Quantity: product.Articles[1].Quantity,
			},
		}
		fx.articlesRepo.EXPECT().RemoveArticles(fx.ctx, articles).Return(nil)

		err := fx.RemoveProduct(fx.ctx, productID, 1)

		require.NoError(t, err)
	})

	t.Run("should remove N product", func(t *testing.T) {
		fx := newFixture(t)

		quantity := int32(testhelpers.RandomIntRange(2, 10))
		fx.productsRepo.EXPECT().GetProduct(fx.ctx, productID).Return(product, nil)
		articles := []models.ProductArticle{
			{
				ID:       product.Articles[0].ID,
				Quantity: product.Articles[0].Quantity * quantity,
			},
			{
				ID:       product.Articles[1].ID,
				Quantity: product.Articles[1].Quantity * quantity,
			},
		}
		fx.articlesRepo.EXPECT().RemoveArticles(fx.ctx, articles).Return(nil)

		err := fx.RemoveProduct(fx.ctx, productID, quantity)

		require.NoError(t, err)
	})
}

type fixture struct {
	Service

	t            *testing.T
	ctx          context.Context
	articlesRepo *mockArticlesRepo.MockRepository
	productsRepo *mockProductsRepo.MockRepository
}

func newFixture(t *testing.T) *fixture {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	fx := &fixture{
		t:            t,
		ctx:          ctx,
		articlesRepo: mockArticlesRepo.NewMockRepository(ctrl),
		productsRepo: mockProductsRepo.NewMockRepository(ctrl),
	}
	fx.Service = NewService(fx.articlesRepo, fx.productsRepo)
	return fx
}
