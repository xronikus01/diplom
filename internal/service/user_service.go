package service

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"blog-api/internal/model"
	"blog-api/internal/repository"
	"blog-api/pkg/auth"
)

type UserService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string, tokenTTL time.Duration) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (s *UserService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if err := validateUsername(username); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	existingByEmail, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingByEmail != nil {
		return nil, ErrUserAlreadyExists
	}

	existingByUsername, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if err := validateEmail(email); err != nil {
		return "", ErrInvalidCredentials
	}
	if password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := auth.CheckPassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func validateUsername(username string) error {
	if len(username) < 3 {
		return ErrInvalidInput
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidInput
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidInput
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidInput
	}

	var hasLetter bool
	var hasDigit bool

	for _, r := range password {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return ErrInvalidInput
	}

	return nil
}
