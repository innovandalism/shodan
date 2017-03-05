package redis

import (
	"gopkg.in/redis.v5"
	"time"
	"fmt"
)

type KVS struct{
	client *redis.Client
	defaultCacheTimer time.Duration
	prefix string
}

// Initialize the redis client for use, call this from the main goroutine only
func (r *KVS) Init(options *redis.Options) (error) {
	r.client = redis.NewClient(options)
	r.prefix = "SHODAN"
	return r.client.Ping().Err()
}

func (r *KVS) formatPrefix(key string) string {
	return fmt.Sprintf("%s_%s", r.prefix, key)
}

func (r *KVS) SetWithDefaultTTL(key string, value interface{}) error {
	return r.client.Set(r.formatPrefix(key), value, r.defaultCacheTimer).Err()
}

func (r *KVS) SetPerm(key string, value interface{}) error {
	return r.client.Set(r.formatPrefix(key), value, time.Duration(0)).Err()
}

func (r *KVS) GetString(key string) (string, error) {
	return r.client.Get(r.formatPrefix(key)).Result()
}

func (r *KVS) HasKey(key string) (bool, error) {
	return r.client.Exists(r.formatPrefix(key)).Result()
}

func (r *KVS) Clear(key string) (int64, error) {
	return r.client.Del(r.formatPrefix(key)).Result()
}