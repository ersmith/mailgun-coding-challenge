package test

import (
	"context"
	"testing"

	"github.com/ersmith/mailgun-coding-challenge/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Test database config
var DatabaseConfig = config.DbConfig{
	Username:    "postgres",
	Password:    "postgres",
	Host:        "localhost",
	Port:        "25432",
	Name:        "mailgun_test",
	MinPoolSize: 1,
	MaxPoolSize: 1,
}

// Creates a connection pool for testing
func CreateTestPgxPool(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), DatabaseConfig.ConnnectionUrl())
	if err != nil {
		t.Fatal(err)
	}

	return pool
}
