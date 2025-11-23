package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/govnocods/RedChat/internal/repository"
	"github.com/govnocods/RedChat/models"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceI interface {
	RegisterUser(username, password string) (*models.User, error)
	Authenticate(username, password string) (*models.User, error)
	ValidateUser(username, password string) error
}

type UserService struct {
	repo repository.UserRepositoryI
}

func NewUserService(repo repository.UserRepositoryI) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(username, password string) (*models.User, error) {
	if existingUser, err := s.repo.GetUser(username); existingUser != nil && err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, fmt.Errorf("repository create user error: %w", err)
	}

	return s.repo.GetUser(username)
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {

	user, err := s.repo.GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) ValidateUser(username, password string) error {
	trimmedUsername := strings.TrimSpace(username)
	if trimmedUsername == "" {
		return errors.New("имя пользователя не может быть пустым")
	}

	if strings.Contains(username, " ") {
		return errors.New("имя пользователя не может содержать пробелы")
	}

	for _, char := range username {
		isLetter := (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
		isDigit := char >= '0' && char <= '9'

		if !isLetter && !isDigit {
			return errors.New("имя пользователя может содержать только английские буквы и цифры")
		}
	}

	usernameLen := utf8.RuneCountInString(username)
	if usernameLen < 4 {
		return errors.New("имя пользователя должно быть длиной не менее 4 символов")
	}

	if len(password) <= 8 {
		return errors.New("пароль должен быть длиннее 8 символов")
	}

	return nil
}
