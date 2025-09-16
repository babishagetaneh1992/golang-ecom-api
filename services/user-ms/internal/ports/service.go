package ports

import "github.com/babishagetaneh1992/ecom-api/services/user-ms/internal/domain"

//import "user-microservice/internal/domain"

// Inbound port (use cases)
type UserService interface {
	Register(user *domain.User) (*domain.User, error)
	GetUser(id string) (*domain.User, error)
	ListUsers() ([]domain.User, error)
	Exists(id string) (bool, error)
	Authenticate(email, password string)(*domain.User, error)
}

