package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("valid db env vars", func(t *testing.T) {
		setDBEnvVars()
		cfg, err := Load()
		assert.NoError(t, err)

		assert.NotNil(t, cfg)
		assert.NotNil(t, cfg.DatabaseConfig)
		assert.Equal(t, "localhost", cfg.DatabaseConfig.Host)
	})

	t.Run("missing db env vars", func(t *testing.T) {
		clearDBEnvVars()
		_, err := Load()
		assert.Error(t, err)
	})

}

func setDBEnvVars() {
	os.Setenv(DBHostKey, "localhost")
	os.Setenv(DBPortKey, "5432")
	os.Setenv(DBNameKey, "dbname")
	os.Setenv(DBUserKey, "user")
	os.Setenv(DBPasswordKey, "password")
}

func clearDBEnvVars() {
	os.Unsetenv(DBHostKey)
	os.Unsetenv(DBPortKey)
	os.Unsetenv(DBNameKey)
	os.Unsetenv(DBUserKey)
	os.Unsetenv(DBPasswordKey)
}
