package repository

import (
	"gorm.io/gorm"
)
func PreloadScope(keys ...string) func(*gorm.DB) *gorm.DB {
  return func(db *gorm.DB) *gorm.DB {
    for _, k := range keys {
      db = db.Preload(k)
    }
    return db
  }
}

// Other scopes here