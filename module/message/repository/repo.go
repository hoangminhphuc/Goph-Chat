package repository

import (
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

type sqlRepo struct {
	db *gorm.DB
}

func NewSQLRepo(db *gorm.DB) *sqlRepo {
	return &sqlRepo{db: db}
}

type redisRepo struct {
	rdb *redis.Client
}

func NewRedisRepo(rdb *redis.Client) *redisRepo {
	return &redisRepo{rdb: rdb}
}