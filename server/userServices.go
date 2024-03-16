package main

import (
	"context"

	"github.com/go-redis/redis/v8" // Import the Redis client package
)

// UserService defines the interface for user-related operations
type UserService interface {
	UpdateScore(entityID string, score int) error
}

// RedisUserService implements the UserService interface using Redis
type RedisUserService struct {
	redisClient *redis.Client // RedisClient instance
}

// NewRedisUserService creates a new RedisUserService instance
func NewRedisUserService(redisClient *redis.Client) *RedisUserService {
	return &RedisUserService{redisClient: redisClient}
}

// UpdateScore updates the score for a user in Redis
func (us *RedisUserService) UpdateScore(entityID string, score int) error {
	// Use HSET command to set the score for the user with the given entity ID
	_, err := us.redisClient.HSet(context.Background(), "users", entityID, score).Result()
	if err != nil {
		return err
	}

	return nil
}
