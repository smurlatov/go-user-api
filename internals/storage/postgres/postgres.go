package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"user-api-service/internals/config"
	"user-api-service/internals/models"
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

	//TODO move it to dockerFile
	instalExtentioneQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`
	_, err = db.Exec(context.Background(), instalExtentioneQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			firstname TEXT NOT NULL,
			lastname TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err = db.Exec(context.Background(), createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(user models.User) (string, error) {
	const op = "storage.postgres.SaveUser"

	createUserQuery := `
		INSERT INTO users (firstname,lastname,email,age)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	_, err := s.db.Prepare(context.Background(), "saveUser", createUserQuery)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var id string
	err = s.db.QueryRow(context.Background(), "saveUser", user.FirstName, user.LastName, user.Email, user.Age).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(id string) (*models.User, error) {
	const op = "storage.postgres.GetUser"
	getUserQuery := `SELECT firstname,lastname,email,age FROM users WHERE id = $1`
	_, err := s.db.Prepare(context.Background(), "getUser", getUserQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var user models.User
	err = s.db.QueryRow(context.Background(), "getUser", id).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (s *Storage) UpdateUser(user models.User, id string) error {
	const op = "storage.postgres.UpdateUser"
	updateUserQuery := `UPDATE users
		SET firstname = $1, lastname = $2, email = $3, age = $4
		WHERE id = $5`
	_, err := s.db.Prepare(context.Background(), "updateUser", updateUserQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(context.Background(), "updateUser", user.FirstName, user.LastName, user.Email, user.Age, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"
	err := s.db.Close(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func getDSNFromConfig(cfg config.Database) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}
