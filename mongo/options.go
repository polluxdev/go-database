package mongo

import "time"

// Option -.
type Option func(*Mongo)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Mongo) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Mongo) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Mongo) {
		c.connTimeout = timeout
	}
}

// DBName -.
func DBName(name string) Option {
	return func(c *Mongo) {
		c.dbName = name
	}
}
