package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	u "github.com/Wave-95/boards/server/internal/api/user"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrBadLogin = errors.New("Could not login with provided credentials")
)

type Service interface {
	Login(ctx context.Context, email, password string) (token string, err error)
}

type service struct {
	userRepo      u.Repository
	jwtSigningKey string
	expiration    int
}

func (s *service) Login(ctx context.Context, email, password string) (token string, err error) {
	user, err := s.userRepo.GetUserByLogin(ctx, email, password)
	if err != nil {
		if errors.Is(err, u.ErrUserDoesNotExist) {
			return "", ErrBadLogin
		}
		return "", fmt.Errorf("service: failed to login: %w", err)
	}
	return s.generateToken(user.Id.String())
}

func NewService(userRepo u.Repository, jwtSigningKey string, expiration int) Service {
	return &service{
		userRepo:      userRepo,
		jwtSigningKey: jwtSigningKey,
		expiration:    expiration,
	}
}

func (s *service) generateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Duration(s.expiration) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSigningKey))
}
