package products

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"warehouse/internal/models"
	"warehouse/internal/testhelpers"
)

func TestImpl_GetProducts(t *testing.T) {
	t.Run("should return empty list", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		items, err := fx.GetProducts(fx.ctx)

		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("should return items", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		p1 := fx.createProduct()
		p2 := fx.createProduct()
		p3 := fx.createProduct()

		items, err := fx.GetProducts(fx.ctx)

		require.NoError(t, err)
		expected := []models.Product{p1, p2, p3}
		assert.Equal(t, expected, items)
	})
}

func TestImpl_GetProduct(t *testing.T) {
	t.Run("should return ErrNotFound", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		fx.createProduct()

		item, err := fx.GetProduct(fx.ctx, testhelpers.RandomInt32())

		require.Equal(t, ErrNotFound, err)
		assert.Empty(t, item)
	})

	t.Run("should get existing product", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		product := fx.createProduct()

		item, err := fx.GetProduct(fx.ctx, product.ID)

		require.NoError(t, err)
		assert.Equal(t, product, item)
	})
}

type fixture struct {
	Repository

	t   *testing.T
	ctx context.Context
	db  *pgxpool.Pool
}

func newFixture(t *testing.T) *fixture {
	ctx := context.Background()
	db := testhelpers.NewDB(t)

	_, err := db.Exec(ctx, "TRUNCATE TABLE products")
	require.NoError(t, err)

	return &fixture{
		t:          t,
		ctx:        ctx,
		db:         db,
		Repository: NewRepository(db),
	}
}

func (fx *fixture) Finish() {
	fx.db.Close()
}

func (fx *fixture) createProduct() models.Product {
	item := models.Product{
		Name:  testhelpers.RandomString(),
		Price: float32(testhelpers.RandomInt()),
		Articles: []models.ProductArticle{
			{
				ID:       testhelpers.RandomInt32(),
				Quantity: testhelpers.RandomInt32(),
			},
		},
	}

	const query = `INSERT INTO products (name, price, articles) VALUES ($1, $2, $3) RETURNING id`
	err := fx.db.QueryRow(fx.ctx, query, item.Name, item.Price, item.Articles).Scan(&item.ID)
	require.NoError(fx.t, err)
	return item
}
