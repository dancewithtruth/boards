package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/Wave-95/boards/backend-core/pkg/security"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
)

var (
	errBadLogin = errors.New("Incorrect email or password.")
)

// Service defines the authentication service interface.
type Service interface {
	Login(ctx context.Context, input LoginInput) (token string, err error)
}

type service struct {
	userRepo   user.Repository
	jwtService jwt.Service
	validator  validator.Validate
}

// NewService creates a new instance of the authentication service.
func NewService(userRepo user.Repository, jwtService jwt.Service, validator validator.Validate) Service {
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
// If there is any other error, it returns a wrapped error.
func (s *service) Login(ctx context.Context, input LoginInput) (token string, err error) {
	if err := s.validator.Struct(input); err != nil {
		return "", err
	}
	retrievedUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", errBadLogin
		}
		return "", fmt.Errorf("service: failed to get user by email: %w", err)
	}
	if ok := security.CheckPasswordHash(input.Password, *retrievedUser.Password); ok == false {
		return "", errBadLogin
	}
	return s.jwtService.GenerateToken(retrievedUser.ID.String())
}
