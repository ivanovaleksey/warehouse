package articles

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"warehouse/internal/models"
	"warehouse/internal/testhelpers"
)

func TestImpl_GetArticles(t *testing.T) {
	t.Run("should return empty list", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		items, err := fx.GetArticles(fx.ctx)

		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("should return items", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		a1 := fx.createArticle(models.Article{})
		a2 := fx.createArticle(models.Article{})
		a3 := fx.createArticle(models.Article{})

		items, err := fx.GetArticles(fx.ctx)

		require.NoError(t, err)
		expected := []models.Article{a1, a2, a3}
		assert.Equal(t, expected, items)
	})
}

func TestImpl_GetArticle(t *testing.T) {
	t.Run("should return ErrNotFound", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		art := models.Article{
			Name:  testhelpers.RandomString(),
			Stock: int32(testhelpers.RandomIntRange(1, 100)),
		}
		art = fx.createArticle(art)

		item, err := fx.GetArticle(fx.ctx, testhelpers.RandomInt32())

		require.Equal(t, ErrNotFound, err)
		assert.Empty(t, item)
	})

	t.Run("should get existing product", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		art := models.Article{
			Name:  testhelpers.RandomString(),
			Stock: int32(testhelpers.RandomIntRange(1, 100)),
		}
		art = fx.createArticle(art)

		item, err := fx.GetArticle(fx.ctx, art.ID)

		require.NoError(t, err)
		assert.Equal(t, art, item)
	})
}

func TestImpl_RemoveArticles(t *testing.T) {
	t.Run("should remove items", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		art1 := fx.createArticle(models.Article{Stock: 10})
		art2 := fx.createArticle(models.Article{Stock: 10})
		art3 := fx.createArticle(models.Article{Stock: 10})
		art4 := fx.createArticle(models.Article{Stock: 10})

		toRemove := []models.ProductArticle{
			{
				ID:       art2.ID,
				Quantity: 10,
			},
			{
				ID:       art3.ID,
				Quantity: 15,
			},
			{
				ID:       art4.ID,
				Quantity: 4,
			},
			{
				ID:       testhelpers.RandomInt32(),
				Quantity: 1,
			},
		}
		err := fx.RemoveArticles(fx.ctx, toRemove)
		require.NoError(t, err)

		art1, err = fx.GetArticle(fx.ctx, art1.ID)
		require.NoError(t, err)
		assert.EqualValues(t, 10, art1.Stock)

		art2, err = fx.GetArticle(fx.ctx, art2.ID)
		require.NoError(t, err)
		assert.EqualValues(t, 0, art2.Stock)

		art3, err = fx.GetArticle(fx.ctx, art3.ID)
		require.NoError(t, err)
		assert.EqualValues(t, -5, art3.Stock)

		art4, err = fx.GetArticle(fx.ctx, art4.ID)
		require.NoError(t, err)
		assert.EqualValues(t, 6, art4.Stock)
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

	_, err := db.Exec(ctx, "TRUNCATE TABLE articles")
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

func (fx *fixture) createArticle(item models.Article) models.Article {
	if (item == models.Article{}) {
		item.Name = testhelpers.RandomString()
		item.Stock = int32(testhelpers.RandomIntRange(1, 100))
	}
	const query = `INSERT INTO articles (name, stock) VALUES ($1, $2) RETURNING id`
	err := fx.db.QueryRow(fx.ctx, query, item.Name, item.Stock).Scan(&item.ID)
	require.NoError(fx.t, err)
	return item
}
