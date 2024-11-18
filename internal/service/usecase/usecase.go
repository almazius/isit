package usecase

import "isit/internal/auth/repository/postgtres"

type Service struct {
	repo postgtres.Repository
}
