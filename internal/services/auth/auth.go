package services

import (
	"errors"
	"log/slog"

	domain "github.com/RamzittoRamzotti/gotochat.git/internal/domain/models"
	"github.com/RamzittoRamzotti/gotochat.git/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidCredentials = errors.New("Invalid password")
	ErrUserNotFound       = errors.New("User not found")
	ErrUserAlreadyExists  = errors.New("User already exists")
)

type AuthService struct {
	userProvider UserProvider
	logger       *slog.Logger
}

type UserProvider interface {
	GetUser(username string) (*domain.User, error)
	CreateUser(username, password string) (*domain.UserCreate, error)
}

func NewAuthService(userProvider UserProvider, logger *slog.Logger) *AuthService {
	return &AuthService{
		userProvider: userProvider,
		logger:       logger,
	}
}

func (s *AuthService) Login(username, password string) (string, error) {
	const op = "AuthService.Login"
	user, err := s.userProvider.GetUser(username)
	if err != nil {
		s.logger.Error("Error occurred while fetching user", slog.String("error", err.Error()))
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		s.logger.Info("Invalid password ", slog.String("username", username))
		return "", errInvalidCredentials
	}
	jwtToken, err := jwt.GenerateToken(user.Username)
	if err != nil {
		s.logger.Error("Error occurred while generating JWT token", slog.String("error", err.Error()))
		return "", err
	}
	return jwtToken, nil
}

func (s *AuthService) Register(username, password string) error {
	const op = "AuthService.Register"
	_, err := s.userProvider.GetUser(username)
	if err == nil {
		s.logger.Info("User already exists", slog.String("username", username))
		return errors.New("User already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error occurred while hashing password", slog.String("error", err.Error()))
		return err
	}
	_, err = s.userProvider.CreateUser(username, string(hashedPassword))
	if err != nil {
		s.logger.Error("Error occurred while creating user", slog.String("error", err.Error()))
		return err
	}
	return nil
}
