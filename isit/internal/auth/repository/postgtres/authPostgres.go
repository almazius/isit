package postgtres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"isit/isit/internal/auth/models"
	"log/slog"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) CreateUser(ctx context.Context, user *models.RegisterParams, hashedPassword string) (int, error) {
	const query = `
INSERT INTO users.user (name, surname, login, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id;
`
	var userId int

	err := repo.db.GetContext(ctx, &userId, query, user.Name, user.Surname, user.Login, hashedPassword)
	if err != nil {
		slog.Error("failed execute query for create user", "error", err)
		return 0, fmt.Errorf("failed exeute query for create user: %w", err)
	}

	return userId, nil
}

func (repo *AuthRepository) GetUserPassword(ctx context.Context, login string) (string, error) {
	const query = `
SELECT password_hash FROM users.user WHERE login = $1 AND is_baned IS FALSE;
`
	var password string
	err := repo.db.GetContext(ctx, &password, query, login)
	if err != nil {
		slog.Error("failed execute query for validate user", "error", err)
		return "", fmt.Errorf("failed exeute query for validate user: %w", err)
	}

	return password, nil
}
