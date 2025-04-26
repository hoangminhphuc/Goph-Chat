package repository

import "gorm.io/gorm"

type sqlRepo struct {
	db *gorm.DB
}

func NewSQLRepo(db *gorm.DB) *sqlRepo {
	return &sqlRepo{db: db}
}