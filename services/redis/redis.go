package redis

import (
	"gopkg.in/redis.v5"
	"github.com/innovandalism/shodan/util"
	"time"
	"fmt"
)

type ShodanRedis struct{
	client *redis.Client
	defaultCacheTimer time.Duration
	prefix string
}

// Initialize the redis client for use, call this from the main goroutine only
func (r *ShodanRedis) Init(options *redis.Options) {
	r.client = redis.NewClient(options)
	r.prefix = "SHODAN"
	err := r.client.Ping().Err()
	util.PanicOnError(err)
}

func (r *ShodanRedis) formatPrefix(key string) string {
	return fmt.Sprintf("%s_%s", r.prefix, key)
}

func (r *ShodanRedis) SetWithDefaultTTL(key string, value interface{}) error {
	return r.client.Set(r.formatPrefix(key), value, r.defaultCacheTimer).Err()
}

func (r *ShodanRedis) SetPerm(key string, value interface{}) error {
	return r.client.Set(r.formatPrefix(key), value, time.Duration(0)).Err()
}

func (r *ShodanRedis) GetString(key string) (string, error) {
	return r.client.Get(r.formatPrefix(key)).Result()
}

func (r *ShodanRedis) HasKey(key string) (bool, error) {
	return r.client.Exists(r.formatPrefix(key)).Result()
}

func (r *ShodanRedis) Clear(key string) (int64, error) {
	return r.client.Del(r.formatPrefix(key)).Result()
}