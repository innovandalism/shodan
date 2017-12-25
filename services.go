package shodan

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	redis "gopkg.in/redis.v5"

	// posthgres driver for database/sql
	_ "github.com/lib/pq"
)

var (
	// DiscordAuthTokenFormat describes the format used for the Authentication header used by discordgo
	DiscordAuthTokenFormat = "Bot %s"
)

// InitCommandStack creates a new CommandStack
func InitCommandStack() (*CommandStack, error) {
	cs := CommandStack{}
	return &cs, nil
}

// InitDiscord sets up a discordgo session (but does not open a connection to a websocket)
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
	log.Printf("Discord Bot %s is connected to\n", discord.State.User.Username)
	userGuilds, err := discord.UserGuilds(100, "", "")
	if err != nil {
		return nil, err
	}
	for _, userGuild := range userGuilds {
		log.Printf(" - [%s] %s\n", userGuild.ID, userGuild.Name)
	}
	return discord, nil
}

// InitHTTP creates a new muxer and starts the HTTP listener
//
// TODO: Reimplement support for disabling the listener
func InitHTTP(addr string) (*mux.Router, error) {
	mux := mux.NewRouter()
	go http.ListenAndServe(addr, mux)
	return mux, nil
}

// InitPostgres initializes the postgres service
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

// InitRedis initializes the Redis KVS
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

// KVS is the inteface for Key-Value-Stores
type KVS interface {
	Set(string, string) error
	SetWithDefaultTTL(string, string) error
	Get(string) (string, error)
	HasKey(key string) (bool, error)
	Clear(string) error
}

// RedisKVS is the redis-implementation of KVS
type RedisKVS struct {
	client            *redis.Client
	defaultCacheTimer time.Duration
	prefix            string
}

func (r *RedisKVS) formatPrefix(key string) string {
	return fmt.Sprintf("%s_%s", r.prefix, key)
}

// SetWithDefaultTTL sets a key with the default expiry timer
func (r *RedisKVS) SetWithDefaultTTL(key string, value string) error {
	return r.client.Set(r.formatPrefix(key), value, r.defaultCacheTimer).Err()
}

// Set sets a key in the KVS permanently
func (r *RedisKVS) Set(key string, value string) error {
	return r.client.Set(r.formatPrefix(key), value, time.Duration(0)).Err()
}

// Get returns a key from the KVS, or an error if the key doesn't exist
func (r *RedisKVS) Get(key string) (string, error) {
	return r.client.Get(r.formatPrefix(key)).Result()
}

// HasKey returns a boolean indicating if the key exists
func (r *RedisKVS) HasKey(key string) (bool, error) {
	return r.client.Exists(r.formatPrefix(key)).Result()
}

// Clear removes a key from the KVS
func (r *RedisKVS) Clear(key string) error {
	return r.client.Del(r.formatPrefix(key)).Err()
}
