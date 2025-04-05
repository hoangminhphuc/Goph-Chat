package models

import "time"

type BaseModel struct {
	ID        int        `json:"-" gorm:"column:id;"`
	CreatedAt *time.Time `json:"-" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at;"`
}