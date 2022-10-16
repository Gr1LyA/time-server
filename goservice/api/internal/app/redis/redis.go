package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type Redis struct {
	*redis.Client
}

// Store - set data into Redis
func (s *Redis) Store(key string, value any) (err error) {
	err = s.Set(context.TODO(), key, value, 0).Err()
	return
}

// Load - Get data from Redis
func (s *Redis) Load(key string) (string, error) {
	content, err := s.Get(context.TODO(), key).Result()
	if err != nil {
		return "", err
	}
	return content, nil
}

// NewRedisClient return store.Redis
func NewRedisClient(addr, passwd string, db int) (*Redis, error) {
	s:= &Redis{redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       0,
	})}

	// check connect
	err := s.Ping(context.TODO()).Err()
	if err != nil {
		return nil, err
	}

	return s, nil
}
