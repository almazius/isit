package usecase

import "isit/isit/internal/auth/repository/postgtres"

type Service struct {
	repo postgtres.Repository
}
