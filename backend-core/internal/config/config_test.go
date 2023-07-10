package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("valid db env vars", func(t *testing.T) {
		cfg, err := Load("../../.env")
		assert.NoError(t, err)

		assert.NotNil(t, cfg)
		assert.NotNil(t, cfg.DB)
		assert.NotEmpty(t, cfg.DB.Host)
	})
}
