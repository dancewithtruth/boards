package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/server/internal/jwt"
	u "github.com/Wave-95/boards/server/internal/user"
)

var (
	ErrBadLogin = errors.New("Could not login with provided credentials")
)

type Service interface {
	Login(ctx context.Context, input LoginInput) (token string, err error)
}

type service struct {
	userRepo   u.Repository
	jwtService jwt.Service
}

func (s *service) Login(ctx context.Context, input LoginInput) (token string, err error) {
	user, err := s.userRepo.GetUserByLogin(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, u.ErrUserDoesNotExist) {
			return "", ErrBadLogin
		}
		return "", fmt.Errorf("service: failed to login: %w", err)
	}
	return s.jwtService.GenerateToken(user.Id.String())
}

func NewService(userRepo u.Repository, jwtService jwt.Service) Service {
	return &service{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}
