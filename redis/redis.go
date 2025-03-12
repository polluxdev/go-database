package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Redis -.
type Redis struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	connPassword string

	DB *redis.Client
}

// New -.
func New(dsn string, opts ...Option) (*Redis, error) {
	rd := &Redis{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(rd)
	}

	var (
		db  *redis.Client
		err error
	)

	for rd.connAttempts > 0 {
		db = redis.NewClient(&redis.Options{
			Addr:           dsn,
			Password:       rd.connPassword,
			DB:             0,
			PoolSize:       rd.maxPoolSize,
			MaxActiveConns: rd.maxPoolSize,
			MaxIdleConns:   rd.maxPoolSize / 2,
			DialTimeout:    rd.connTimeout,
		})
		_, err = db.Ping(context.Background()).Result()
		if err == nil {
			break
		}

		log.Printf("Redis is trying to connect, attempts left: %d", rd.connAttempts)

		time.Sleep(rd.connTimeout)

		rd.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis - New - connAttempts == 0: %w", err)
	}

	log.Println("Redis is connected successfully")

	_, err = db.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis - New - Ping: %w", err)
	}

	rd.DB = db

	return rd, nil
}

// Close -.
func (m *Redis) Close() {
	if m.DB != nil {
		m.DB.WithTimeout(m.connTimeout).Close()
	}
}
