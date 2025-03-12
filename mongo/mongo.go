package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Mongo -.
type Mongo struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	dbName       string

	client *mongo.Client
	DB     *mongo.Database
}

// New -.
func New(dsn string, opts ...Option) (*Mongo, error) {
	mg := &Mongo{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(mg)
	}

	var (
		db  *mongo.Client
		err error
	)

	for mg.connAttempts > 0 {
		db, err = mongo.Connect(
			options.MergeClientOptions(
				options.Client().ApplyURI(dsn),
				options.Client().SetMaxPoolSize(uint64(mg.maxPoolSize)),
				options.Client().SetMaxConnIdleTime(time.Duration(mg.maxPoolSize)),
				options.Client().SetConnectTimeout(mg.connTimeout),
			),
		)
		if err == nil {
			break
		}

		log.Printf("Mongo is trying to connect, attempts left: %d", mg.connAttempts)

		time.Sleep(mg.connTimeout)

		mg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("mongo - New - connAttempts == 0: %w", err)
	}

	log.Println("Mongo is connected successfully")

	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("mongo - New - Ping: %w", err)
	}

	mg.client = db
	mg.DB = db.Database(mg.dbName)

	return mg, nil
}

// Close -.
func (m *Mongo) Close() {
	if m.DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), m.connTimeout)
		defer cancel()

		m.client.Disconnect(ctx)
	}
}
