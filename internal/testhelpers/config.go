package testhelpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"warehouse/internal/config"
)

func NewConfig(t *testing.T) config.Config {
	return config.NewConfigFromFile(findConfigFile(t, "config/config.test.yaml"))
}

func findConfigFile(t *testing.T, f string) string {
	path := "."
	for {
		tryPath := filepath.Join(path, f)
		var err error
		if _, err = os.Stat(tryPath); err == nil {
			return tryPath
		}
		require.True(t, os.IsNotExist(err))
		path = filepath.Join("..", path)
	}
}
