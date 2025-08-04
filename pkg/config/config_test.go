package config

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = godotenv.Load("../../.env")
}

func TestConfigValidity(t *testing.T) {
	cfg, err := GetConfig()
	assert.NoError(t, err)

	t.Logf("DatabaseFile %v", cfg.DatabaseFile)
}
