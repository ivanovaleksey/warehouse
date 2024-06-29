package testhelpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"warehouse/internal/db"
)

func NewDB(t *testing.T) *pgxpool.Pool {
	appCfg := NewConfig(t)
	cfg, err := db.ParseConfig(appCfg)
	require.NoError(t, err)
	pool, err := pgxpool.New(context.Background(), cfg.DSN())
	require.NoError(t, err)
	return pool
}
