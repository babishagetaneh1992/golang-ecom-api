package application

import (
	"errors"
	"fmt"
	"time"

	//"user-microservice/internal/domain"
	//"user-microservice/internal/ports"

	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/domain"
	"github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/ports"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImplement struct {
	repo ports.UserRepository
	validate *validator.Validate
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return  &UserServiceImplement{
		repo: repo,
		validate: validator.New(),
	}
}


func (s *UserServiceImplement) Register(user *domain.User) (*domain.User, error) {
	// validate struct
	if err := s.validate.Struct(user); err != nil {
		return nil, err
	}

	// check if email already exists
	existing, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	user.CreateAt = time.Now()

	// insert
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}


func (s *UserServiceImplement) GetUser(id string) (*domain.User, error) {
   return  s.repo.GetById(id)
}

func (s *UserServiceImplement) ListUsers() ([]domain.User, error) {
	return  s.repo.GetAll()
}
func (s *UserServiceImplement) Exists(id string) (bool, error)  {
    return  s.repo.Exists(id)
}

func (s *UserServiceImplement) Authenticate(email, password string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return  nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return  nil, errors.New("invalid email or password")
	}

	return user, nil
}
