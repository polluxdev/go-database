package postgres

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *gorm.DB
}

// New -.
func New(dsn string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	var (
		db  *gorm.DB
		err error
	)

	for pg.connAttempts > 0 {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Info),
			TranslateError: true,
		})
		if err == nil {
			break
		}

		log.Printf("PostgreSQL is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - connAttempts == 0: %w", err)
	}

	log.Println("PostgreSQL is connected successfully")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres - New - db.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(pg.maxPoolSize)
	sqlDB.SetMaxIdleConns(pg.maxPoolSize / 2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	pg.DB = db

	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.DB != nil {
		db, _ := p.DB.DB()
		db.Close()
	}
}
