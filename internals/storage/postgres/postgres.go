package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"user-api-service/internals/config"
)

type Storage struct {
	db *pgx.Conn
}

func Connect(cfg config.Database) (*Storage, error) {
	const op = "storage.postgres.Connect"
	db, err := pgx.Connect(context.Background(), getDSNFromConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			firstname TEXT NOT NULL,
			lastname TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER NOT NULL,
			created TIMESTAMP NOT NULL
		)
	`
	_, err = db.Exec(context.Background(), createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func getDSNFromConfig(cfg config.Database) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}
