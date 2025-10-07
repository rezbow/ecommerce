package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rezbow/ecommerce/internal/app/models"
)

type ICartRepo interface {
	Get(string) (*models.Cart, error)
	Save(string, *models.Cart, time.Duration) error
	Delete(string) error
}

type CartRepoRedis struct {
	client *redis.Client
	ctx    context.Context
}

func NewCartRepoRedis(client *redis.Client) *CartRepoRedis {
	return &CartRepoRedis{
		client: client,
	}
}

func (repo *CartRepoRedis) Get(key string) (*models.Cart, error) {
	value, err := repo.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, ErrInternal
	}

	var cart models.Cart
	if err := json.Unmarshal([]byte(value), &cart); err != nil {
		return nil, ErrInternal
	}
	return &cart, nil
}

func (repo *CartRepoRedis) Save(key string, cart *models.Cart, exp time.Duration) error {
	value, err := json.Marshal(cart)
	if err != nil {
		return ErrInternal
	}
	return repo.client.Set(context.Background(), key, string(value), exp).Err()
}

func (repo *CartRepoRedis) Delete(key string) error {
	result, err := repo.client.Del(context.Background(), key).Result()
	if err != nil {
		return ErrInternal
	}
	if result == 0 {
		return ErrRecordNotFound
	}
	return nil
}
