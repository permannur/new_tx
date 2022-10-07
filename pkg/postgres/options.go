package postgres

import "time"

type Option func(postgres *Postgres)

func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

func ConnAttempt(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
