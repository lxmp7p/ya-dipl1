package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationPath(t *testing.T) {
	assert.Equal(t, migrationsPath, "file://migrations")
}
