package auth

import (
	"context"
	"errors"
	"fmt"

	u "github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/jwt"
)

var (
	ErrBadLogin = errors.New("Could not login with provided credentials")
)

type Service interface {
	Login(ctx context.Context, email, password string) (token string, err error)
}

type service struct {
	userRepo   u.Repository
	jwtService jwt.JWTService
}

func (s *service) Login(ctx context.Context, email, password string) (token string, err error) {
	user, err := s.userRepo.GetUserByLogin(ctx, email, password)
	if err != nil {
		if errors.Is(err, u.ErrUserDoesNotExist) {
			return "", ErrBadLogin
		}
		return "", fmt.Errorf("service: failed to login: %w", err)
	}
	return s.jwtService.GenerateToken(user.Id.String())
}

func NewService(userRepo u.Repository, jwtService jwt.JWTService) Service {
	return &service{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}
