package postgtres

import (
	"context"
	"isit/isit/internal/auth/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.RegisterParams, hashedPassword string) (int, error)
	GetUserPassword(ctx context.Context, login string) (string, error)
}
