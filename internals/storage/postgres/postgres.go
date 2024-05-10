package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"user-api-service/internals/config"
	"user-api-service/internals/models"
	"user-api-service/internals/storage/queries"
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
	installExtentionQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`
	_, err = db.Exec(context.Background(), installExtentionQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// create table
	_, err = db.Exec(context.Background(), queries.CreateUserTableQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(user models.User) (string, error) {
	const op = "storage.postgres.SaveUser"
	//TODO move prepare to init func
	_, err := s.db.Prepare(context.Background(), "saveUser", queries.CreateUserQuery)
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

func (s *Storage) GetUser(id string) (models.User, error) {
	const op = "storage.postgres.GetUser"
	_, err := s.db.Prepare(context.Background(), "getUser", queries.SelectUserByIDQuery)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	err = s.db.QueryRow(context.Background(), "getUser", id).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Age)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) UpdateUser(user models.User, id string) error {
	const op = "storage.postgres.UpdateUser"
	_, err := s.db.Prepare(context.Background(), "updateUser", queries.UpdateUserQuery)
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
