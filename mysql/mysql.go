package mysql

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// MySQL -.
type MySQL struct {
	dsn          string
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
}

// MySQLManager -.
type MySQLManager struct {
	mu          sync.Mutex
	connections map[string]*gorm.DB
	configs     map[string]*MySQL
}

// New -.
func New() *MySQLManager {
	return &MySQLManager{
		connections: make(map[string]*gorm.DB),
		configs:     make(map[string]*MySQL),
	}
}

// AddDatabase -.
func (m *MySQLManager) AddDatabase(name, dsn string, opts ...Option) {
	m.mu.Lock()
	defer m.mu.Unlock()

	config := &MySQL{
		dsn:          dsn,
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(config)
	}

	m.configs[name] = config
}

// Connect -.
func (m *MySQLManager) Connect(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if connection already exists
	if _, exists := m.connections[name]; exists {
		return nil // Already connected
	}

	config, exists := m.configs[name]
	if !exists {
		return fmt.Errorf("database config for '%s' not found", name)
	}

	var (
		db       *gorm.DB
		err      error
		attempts = config.connAttempts
	)

	for attempts > 0 {
		db, err = gorm.Open(mysql.Open(config.dsn), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Info),
			TranslateError: true,
		})
		if err == nil {
			break
		}

		log.Printf("MySQL '%s' is trying to connect, attempts left: %d", name, attempts)

		time.Sleep(config.connTimeout)

		attempts--
	}

	if err != nil {
		return fmt.Errorf("mysql - New - connAttempts == 0 for '%s': %w", name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("mysql - New - db.DB for '%s': %w", name, err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(config.maxPoolSize)
	sqlDB.SetMaxIdleConns(config.maxPoolSize / 2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	m.connections[name] = db

	log.Printf("MySQL %s DB is connected successfully\n", name)

	return nil
}

// Get -.
func (m *MySQLManager) GetDB(name string) (*gorm.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("database connection '%s' not found", name)
	}
	return db, nil
}

// CloseAll -.
func (m *MySQLManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, db := range m.connections {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
		delete(m.connections, name)
	}
}
