package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"isit/internal/auth/models"
	"isit/internal/auth/repository/postgtres"
	"isit/internal/auth/repository/redis"
	"log/slog"
	"math/rand"
	"strings"
)

const (
	defaultCost      = 14
	SessionLengthKey = 32
)

type AuthService struct {
	authRepo postgtres.Repository
	cache    *redis.AuthCache
}

func NewAuthService(authRepo postgtres.Repository, cache *redis.AuthCache) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		cache:    cache,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, params *models.RegisterParams) error {
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	userId, err := s.authRepo.CreateUser(ctx, params, hashedPassword)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	_ = userId

	return nil
}

func (s *AuthService) SingIn(ctx context.Context, params *models.AuthParams, userAgent, fingerprint string) (string, error) {
	sessionKey := ""

	passwordHash, err := s.authRepo.GetUserPassword(ctx, params.Login)
	if err != nil {
		slog.Error("failed check user", "error", err)
		return "", err
	}

	if CheckPasswordHash(params.Password, passwordHash) {
		sessionKey = generateSessionKey()
		err = s.cache.CreateNewSession(ctx, fingerprint, userAgent, sessionKey)
		if err != nil {
			slog.Error("failed to create session", "error", err)
			return sessionKey, err
		}
	}

	return sessionKey, nil
}

func (s *AuthService) ValidateSessionKey(ctx context.Context, userAgent, fingerprint, sessionKey string) error {
	key, err := s.cache.ValidateSession(ctx, fingerprint, userAgent)
	if err != nil {
		slog.Error("failed validate session key", "error", err)
		return err
	}

	if sessionKey != key {
		slog.Error("session key is not trusted")
		return fmt.Errorf("session key is not trusted")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}

	return string(pass), nil
}

func generateSessionKey() string {
	dictionary := []rune{
		'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', 'a', 's',
		'd', 'f', 'g', 'h', 'j', 'k', 'l', 'z', 'x', 'c', 'v', 'b',
		'n', 'm', 'Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P',
		'A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'Z', 'X', 'C',
		'V', 'B', 'N', 'M',
		'1', '2', '3', '4', '5',
		'6', '7', '8', '9', '0'}

	builder := strings.Builder{}
	for i := 0; i < SessionLengthKey; i++ {
		builder.WriteRune(dictionary[rand.Int()%len(dictionary)])
	}

	return builder.String()
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
