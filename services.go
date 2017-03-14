package shodan

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gopkg.in/redis.v5"
	"net/http"
	"time"
)

var (
	DiscordAuthTokenFormat = "Bot %s"
)

func InitCommandStack() (*CommandStack, error) {
	cs := CommandStack{}
	return &cs, nil
}

func InitDiscord(token string) (*discordgo.Session, error) {
	t := fmt.Sprintf(DiscordAuthTokenFormat, token)
	discord, err := discordgo.New(t)
	if err != nil {
		return nil, err
	}
	discord.State.User, err = discord.User("@me")
	if err != nil {
		return nil, err
	}
	return discord, nil
}

func InitHTTP(addr string) (*mux.Router, error) {
	mux := mux.NewRouter()
	go http.ListenAndServe(addr, mux)
	return mux, nil
}

// Initialize Postgres Service
func InitPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Connectivity/Auth test
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis(url string) (KVS, error) {
	kvs := RedisKVS{}
	kvs.prefix = "shodan"
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	kvs.client = redis.NewClient(options)
	err = kvs.client.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &kvs, nil
}

type KVS interface {
	Set(string, string) error
	SetWithDefaultTTL(string, string) error
	Get(string) (string, error)
	HasKey(key string) (bool, error)
	Clear(string) error
}

type RedisKVS struct {
	client            *redis.Client
	defaultCacheTimer time.Duration
	prefix            string
}

func (r *RedisKVS) formatPrefix(key string) string {
	return fmt.Sprintf("%s_%s", r.prefix, key)
}

func (r *RedisKVS) SetWithDefaultTTL(key string, value string) error {
	return r.client.Set(r.formatPrefix(key), value, r.defaultCacheTimer).Err()
}

func (r *RedisKVS) Set(key string, value string) error {
	return r.client.Set(r.formatPrefix(key), value, time.Duration(0)).Err()
}

func (r *RedisKVS) Get(key string) (string, error) {
	return r.client.Get(r.formatPrefix(key)).Result()
}

func (r *RedisKVS) HasKey(key string) (bool, error) {
	return r.client.Exists(r.formatPrefix(key)).Result()
}

func (r *RedisKVS) Clear(key string) error {
	return r.client.Del(r.formatPrefix(key)).Err()
}
