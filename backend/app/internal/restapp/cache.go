package restapp

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client  *redis.Client
	context context.Context
}

func NewCache(address string, password string) (*Cache, error) {
	var ctx = context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // use default DB
	})

	_, err := client.Ping(ctx).Result()

	return &Cache{
		client:  client,
		context: ctx,
	}, err
}

func (c *Cache) Get(key string) (string, error) {
	value, err := c.client.Get(c.context, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	return c.client.Set(c.context, key, value, ttl).Err()
}

func (c *Cache) SetInt(key string, value int, ttl time.Duration) error {
	strValue := strconv.Itoa(value)
	return c.Set(key, strValue, ttl)
}

// returns -1 if not found
func (c *Cache) GetInt(key string) (int, error) {
	result, err := c.Get(key)
	if result == "" || err != nil {
		return -1, err
	}
	return strconv.Atoi(result)
}
