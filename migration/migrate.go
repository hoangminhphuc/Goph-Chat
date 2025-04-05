package migration

import (
	"fmt"

	"github.com/hoangminhphuc/goph-chat/module/user/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Debug().AutoMigrate(
		&model.User{},
	); err != nil {
		return fmt.Errorf("failed to migrate %w", err)
	}

	return nil
}
