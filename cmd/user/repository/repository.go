package repository

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserRepository struct {
	Redis    *redis.Client
	Database *gorm.DB
}

func NewUserRepository(redisClient *redis.Client, db *gorm.DB) *UserRepository {
	return &UserRepository{
		Redis:    redisClient,
		Database: db,
	}
}
