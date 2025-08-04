package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Item represents the database model
type Item struct {
	gorm.Model
	Hash     string   `gorm:"type:text;uniqueIndex;not null"` // SHA-256 hash with unique index
	Framing  *float32 `gorm:"type:real"`                      // Nullable
	TitleAI  *string  `gorm:"type:text"`                      // Nullable
	ReasonAI *string  `gorm:"type:text"`                      // Nullable
}

// TableName explicitly sets the table name and ensures schema constraints
func (Item) TableName() string {
	return "items"
}

// Database handles DB operations
type Database struct {
	db *gorm.DB
}

// NewDatabase initializes a new SQLite database
func NewDatabase(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate to create table with constraints
	err = db.AutoMigrate(&Item{})
	if err != nil {
		return nil, err
	}

	// Explicitly ensure unique index on Hash
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_hash ON items(hash)").Error
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// CreateItem inserts a new item, ignores if hash already exists
func (d *Database) CreateItem(item *Item) error {
	return d.db.Exec(
		"INSERT OR IGNORE INTO items (hash, framing, title_ai, reason_ai, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		item.Hash, item.Framing, item.TitleAI, item.ReasonAI, item.CreatedAt, item.UpdatedAt,
	).Error
}

// FindItemByHash retrieves an item by its hash
func (d *Database) FindItemByHash(hash string) (*Item, error) {
	var item Item
	result := d.db.Where("hash = ?", hash).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}
