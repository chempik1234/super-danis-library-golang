package genericport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chempik1234/super-danis-library-golang/v2/pkg/genericports"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisGenericCache - implement genericports.GenericCachePort
type RedisGenericCache[K comparable, V genericports.ObjectWithIdentifier[K]] struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisGenericCache creates a new instance of RedisGenericCache
func NewRedisGenericCache[K comparable, V genericports.ObjectWithIdentifier[K]](addr string, password string, db int, ttlMs int) *RedisGenericCache[K, V] {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisGenericCache[K, V]{client: client, ttl: time.Duration(ttlMs) * time.Millisecond}
}

// GetObjectByID - impl genericports.GenericCachePort.GetObjectByID
func (s *RedisGenericCache[K, V]) GetObjectByID(ctx context.Context, id K) (*V, error) {
	key := generateKey(id)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found
		}
		return nil, err // Other errors
	}

	var obj V
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

// SaveObject - impl genericports.GenericCachePort.SaveObject
func (s *RedisGenericCache[K, V]) SaveObject(ctx context.Context, fullyReadyObject *V) (*V, error) {
	key := generateKey((*fullyReadyObject).GetUniqueIdentifier())
	data, err := json.Marshal(fullyReadyObject)
	if err != nil {
		return nil, err
	}

	if err := s.client.Set(ctx, key, data, s.ttl).Err(); err != nil {
		return nil, err
	}

	return fullyReadyObject, nil
}

// DeleteObject - impl genericports.GenericCachePort.DeleteObject
func (s *RedisGenericCache[K, V]) DeleteObject(ctx context.Context, id K) error {
	key := generateKey(id)
	return s.client.Del(ctx, key).Err()
}

// generateKey generates a Redis key based on the ID
func generateKey[K comparable](id K) string {
	return fmt.Sprintf("gnrc_rds_%v", id)
}
