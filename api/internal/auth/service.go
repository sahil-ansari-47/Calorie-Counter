package auth

import (
	"caloriecounter/api/internal/database/sqlc"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	queries *sqlc.Queries
	tokens  TokenService
}

func NewService(pool *pgxpool.Pool, tokens TokenService) *Service {
	return &Service{queries: sqlc.New(pool), tokens: tokens}
}
func (s *Service) Register(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !strings.Contains(email, "@") || len(password) < 8 {
		return "", ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}, Email: email, PasswordHash: string(hash)})
	if err != nil {
		return "", err
	}
	return s.tokens.Issue(uuid.UUID(user.ID.Bytes))
}
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.queries.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}
	return s.tokens.Issue(uuid.UUID(user.ID.Bytes))
}
