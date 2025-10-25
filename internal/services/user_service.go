package services

import (
	"errors"
	"playspotter/internal/models"
	"playspotter/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateUser(id uuid.UUID, name, password string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if name != "" {
		user.Name = name
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	return s.userRepo.Update(user)
}

func (s *UserService) ListUsers(offset, limit int) ([]models.User, int64, error) {
	return s.userRepo.List(offset, limit)
}

func (s *UserService) UpdateUserRole(id uuid.UUID, role string) error {
	if role != "user" && role != "admin" {
		return errors.New("invalid role")
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	user.Role = role
	return s.userRepo.Update(user)
}
