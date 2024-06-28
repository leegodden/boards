package user

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/leegodden/boards/server/internal/entity"
)

var ErrInvalidEmail = errors.New("invalid email address")

type CreateUserInput struct {
	Name     string
	Email    string `validate:"email"`
	Password string
	IsGuest  bool
}

func (input *CreateUserInput) Validate() error {
	v := validator.New()
	
	if err := v.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Email" {
				return ErrInvalidEmail
			}
		}
	}
	return nil
}

// Any type that wants to satisfy this interface, should have a method named 
// CreateUser which takes CreateUserInput as an argument and returns error."
type Service interface {
	CreateUser(input CreateUserInput) (*User, error)
}

type User struct {
	entity.User
}

// Struct with variable repo of type Repository(defined in repository.go)
type service struct {
	repo Repository
}


// A method that we're associating with the service struct. The portion (s *service) 
// is the receiver and essentially dictates that this function belongs to, and 
// operates upon the service struct. Now our service struct conforms to the Service 
// interface, because it has its own CreateUser method. `s' is he intance of the caller`
func (s *service) CreateUser(input CreateUserInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	id := uuid.New()
	now := time.Now()

	// Create a new instance of User, assign the values from 'input' to the new instance
	// and save the new User instance to the database
	user := entity.User{
		Id:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		IsGuest:   input.IsGuest,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &User{user}, nil
}

// Creates an instance of 'service' and sets the service struct's 'repo' field to the 
// repo that's passed in as a argument. The '& 'before service is used to return a pointer 
// to the new service instance, rather than the instance itself to persist changes 
// across function calls and improve efficiency. Here the Repository interface is the 
// dependency that is being injected into our Service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}