package domain

import (
	"errors"
	"time"
)

type User struct {
	ID       string     `json:"id" bson:"_id,omitempty"`
	Name     string     `json:"name" validate:"required,min=2"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=6"`
	CreateAt time.Time  `json:"created_at" bson:"created_at"`
}


func  NewUser(id, name, email, hashedPassword string) (*User, error) {
	if len(name) < 2 {
		return  nil, errors.New("name must be at least 2 characters")
	}

	if len(hashedPassword) < 6 {
        return  nil, errors.New("password must be 6 characters")
	}

	return &User{
		ID: id,
		Name: name,
		Email: email,
		Password: hashedPassword,
		CreateAt: time.Now(),
	}, nil
}