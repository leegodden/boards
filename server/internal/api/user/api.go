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

var (
	ErrMissingName              = errors.New("missing name")
	ErrInvalidCreateUserRequest = errors.New("invalid create user request")
	ErrInternalServerError      = errors.New("issue creating user")
)

// Holds an instance of the Service interface
type API struct {
	service Service
}

// Takes in a service as argument and sets the service field of the API struct 
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

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode request and validate
	var createUserRequest CreateUserRequest	
	json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err := createUserRequest.Validate(); err != nil {
		response.RespondWithError(w, http.StatusBadRequest, err)
	}

	// Create user
	input := CreateUserInput(createUserRequest)
	user, err := api.service.CreateUser(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write response
	w.WriteHeader(http.StatusCreated)
	createUserResponse := CreateUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		IsGuest:   user.IsGuest,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	json.NewEncoder(w).Encode(createUserResponse)
}