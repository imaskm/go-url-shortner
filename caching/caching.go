package caching

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Cache{
		client: client,
	}
}

func (c *Cache) Write(key, value string) error {

	err := c.client.Set(context.TODO(), key, value, time.Hour*12).Err()

	return err

}

func (c *Cache) Read(key string) (string, error) {

	value, err := c.client.Get(context.TODO(), key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}
