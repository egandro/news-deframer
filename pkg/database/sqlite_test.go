package database

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *Database {
	// Use in-memory SQLite for testing
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestNewDatabase(t *testing.T) {
	db := setupTestDB(t)
	assert.NotNil(t, db, "Database should be initialized")
}

func TestTableExists(t *testing.T) {
	db := setupTestDB(t)

	// Check if the 'items' table exists
	var tableExists bool
	err := db.db.Raw("SELECT count(*) > 0 FROM sqlite_master WHERE type='table' AND name='items'").Scan(&tableExists).Error
	assert.NoError(t, err, "Querying table existence should succeed")
	assert.True(t, tableExists, "Table 'items' should exist after migration")
}

func TestCreateItem(t *testing.T) {
	db := setupTestDB(t)

	// Create test item
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("test")))
	item := &Item{
		Hash:     hash,
		Framing:  new(float32),
		TitleAI:  new(string),
		ReasonAI: new(string),
	}
	*item.Framing = 0.5
	*item.TitleAI = "Test Title"
	*item.ReasonAI = "Test Reason"

	// Insert item
	err := db.CreateItem(item)
	assert.NoError(t, err, "Item creation should succeed")

	// Verify item exists
	found, err := db.FindItemByHash(hash)
	assert.NoError(t, err, "FindItemByHash should succeed")
	assert.NotNil(t, found, "Item should be found")
	assert.Equal(t, hash, found.Hash, "Hash should match")
	assert.Equal(t, *item.Framing, *found.Framing, "Framing should match")
	assert.Equal(t, *item.TitleAI, *found.TitleAI, "TitleAI should match")
	assert.Equal(t, *item.ReasonAI, *found.ReasonAI, "ReasonAI should match")

	// Test creating duplicate hash (should silently ignore)
	duplicate := &Item{
		Hash:     hash,
		Framing:  new(float32),
		TitleAI:  new(string),
		ReasonAI: new(string),
	}
	*duplicate.Framing = 0.7
	*duplicate.TitleAI = "Different Title"
	*duplicate.ReasonAI = "Different Reason"
	err = db.CreateItem(duplicate)
	assert.NoError(t, err, "Creating item with duplicate hash should succeed (ignored)")

	// Verify original item is unchanged
	found, err = db.FindItemByHash(hash)
	assert.NoError(t, err, "FindItemByHash should succeed")
	assert.Equal(t, *item.Framing, *found.Framing, "Original Framing should remain")
	assert.Equal(t, *item.TitleAI, *found.TitleAI, "Original TitleAI should remain")
	assert.Equal(t, *item.ReasonAI, *found.ReasonAI, "Original ReasonAI should remain")
}

func TestFindItemByHash(t *testing.T) {
	db := setupTestDB(t)

	// Create test item
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("test")))
	item := &Item{
		Hash:     hash,
		Framing:  new(float32),
		TitleAI:  new(string),
		ReasonAI: new(string),
	}
	*item.Framing = 0.5
	*item.TitleAI = "Test Title"
	*item.ReasonAI = "Test Reason"

	// Insert item
	err := db.CreateItem(item)
	assert.NoError(t, err, "Item creation should succeed")

	// Test finding item by hash
	found, err := db.FindItemByHash(hash)
	assert.NoError(t, err, "FindItemByHash should succeed")
	assert.NotNil(t, found, "Item should be found")
	assert.Equal(t, hash, found.Hash, "Hash should match")
	assert.Equal(t, *item.Framing, *found.Framing, "Framing should match")
	assert.Equal(t, *item.TitleAI, *found.TitleAI, "TitleAI should match")
	assert.Equal(t, *item.ReasonAI, *found.ReasonAI, "ReasonAI should match")

	// Test non-existent hash
	_, err = db.FindItemByHash("nonexistent")
	assert.Error(t, err, "FindItemByHash should fail for non-existent hash")
}

func TestItemConstraints(t *testing.T) {
	db := setupTestDB(t)

	// Test unique constraint on Hash
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("test")))
	item1 := &Item{
		Hash:     hash,
		Framing:  new(float32),
		TitleAI:  new(string),
		ReasonAI: new(string),
	}
	*item1.Framing = 0.5
	*item1.TitleAI = "Title1"
	*item1.ReasonAI = "Reason1"
	err := db.CreateItem(item1)
	assert.NoError(t, err, "First item creation should succeed")

	item2 := &Item{
		Hash:     hash,
		Framing:  new(float32),
		TitleAI:  new(string),
		ReasonAI: new(string),
	}
	*item2.Framing = 0.7
	*item2.TitleAI = "Title2"
	*item2.ReasonAI = "Reason2"
	err = db.CreateItem(item2)
	assert.NoError(t, err, "Creating item with duplicate Hash should succeed (ignored)")
}
