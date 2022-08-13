package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		c := New()
		assert.Equal(t, "http://localhost:8081", c.RunAddress)
		assert.Equal(t, "", c.DatabaseURI)
	})
}
