package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/leegodden/boards/server/internal/response"
)

/*

api.go defines the HTTP API for the user service, including request/response structures
and handlers. Handles client requests and call the corresponding methods on the service interface,
ensuring separation of concerns.
*/

var (
	ErrMissingName              = errors.New("missing name")
	ErrInvalidCreateUserRequest = errors.New("invalid create user request")
	ErrInternalServerError      = errors.New("issue creating user")
)

// Holds an instance of the Service interface
type API struct {
	service Service
}

// initialize and return a new instance of the `API struct`
func NewAPI(service Service) API {
	return API{service: service}
}

// CreateUserRequest
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsGuest  bool   `json:"is_guest"`
}

// Validates the CreateUserRequest
func (req *CreateUserRequest) Validate() error {
	v := validator.New()
	if err := v.Struct(req); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return ErrMissingName
			}  
		}
		return ErrInvalidCreateUserRequest
	}
	return nil
}

// CreateUserResponse
type CreateUserResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsGuest   bool      `json:"is_guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// processes HTTP requests that aim to create a new user
func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode request and validate
	var createUserRequest CreateUserRequest	
	json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err := createUserRequest.Validate(); err != nil {
		response.WriteWithError(w, http.StatusBadRequest, err)
		return
	}

	// If no errors then create user 
	input := CreateUserInput(createUserRequest)
	user, err := api.service.CreateUser(input)
	if err != nil {
		response.WriteWithError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	// Write the response using "WriteHeader" which sends an HTTP response 
	// header with the provided status code.
	w.WriteHeader(http.StatusCreated)
	createUserResponse := CreateUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		IsGuest:   user.IsGuest,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	response.WriteWithStatus(w, http.StatusCreated, createUserResponse)
}