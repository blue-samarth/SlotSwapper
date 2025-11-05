package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"SlotSwapper/config"
)

func TestHealthCheck_DatabaseConnectivity(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	// Test database connectivity
	sqlDB, err := db.DB()
	assert.NoError(t, err)

	err = sqlDB.Ping()
	assert.NoError(t, err, "Database should be reachable")
}

func TestHealthCheck_DatabaseStats(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	stats := sqlDB.Stats()
	assert.GreaterOrEqual(t, stats.OpenConnections, 0)
	assert.GreaterOrEqual(t, stats.Idle, 0)
}

func TestHealthCheck_ConfigDBAccess(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	// Verify config.GetDB() returns the injected test DB
	configDB := config.GetDB()
	assert.NotNil(t, configDB)

	sqlDB, err := configDB.DB()
	assert.NoError(t, err)

	err = sqlDB.Ping()
	assert.NoError(t, err)
}

func TestHealthCheck_MultipleChecks(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(db)

	// Simulate multiple health checks
	for i := 0; i < 5; i++ {
		sqlDB, err := db.DB()
		assert.NoError(t, err)

		err = sqlDB.Ping()
		assert.NoError(t, err, "Database should remain healthy across multiple checks")
	}
}
