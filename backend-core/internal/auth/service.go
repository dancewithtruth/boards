package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
	u "github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
)

var (
	// ErrBadLogin is an error that is used when a user cannot be found with the given credentials
	ErrBadLogin = errors.New("Could not login with provided credentials")
)

// Service defines the authentication service interface.
type Service interface {
	Login(ctx context.Context, input LoginInput) (token string, err error)
}

type service struct {
	userRepo   u.Repository
	jwtService jwt.Service
	validator  validator.Validate
}

// NewService creates a new instance of the authentication service.
func NewService(userRepo u.Repository, jwtService jwt.Service, validator validator.Validate) Service {
	return &service{
		userRepo:   userRepo,
		jwtService: jwtService,
		validator:  validator,
	}
}

// Login performs the login process using the provided user credentials.
// It retrieves the user from the user repository and generates a JWT token.
// If successful, it returns the generated token.
// If the user cannot be found, it returns ErrBadLogin.
// If there is any other error, it returns a formatted error message.
func (s *service) Login(ctx context.Context, input LoginInput) (token string, err error) {
	if err := s.validator.Struct(input); err != nil {
		return "", err
	}
	user, err := s.userRepo.GetUserByLogin(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, u.ErrUserDoesNotExist) {
			return "", ErrBadLogin
		}
		return "", fmt.Errorf("service: failed to login: %w", err)
	}
	return s.jwtService.GenerateToken(user.ID.String())
}
