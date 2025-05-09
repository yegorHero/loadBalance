package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

const (
	maxRetries = 10
	retryInterval = 2 * time.Second
)

func NewStorage(urlDB string) *Storage {
	var pool *pgxpool.Pool

	for i := 0; i < maxRetries; i++ {

	}

	return *Storage{
		db:
	}
}
