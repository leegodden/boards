package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/leegodden/boards/server/internal/entity"
	"github.com/stretchr/testify/assert"
)

// Testing the service-level logic with a mocked Repository to avoid any database
// dependencies, making tests faster and more reliable

// Call chain ->
// service.CreateUser from service_test.go -> service.CreateUser from service.go ->
// mockRepo.CreateUser from repository_mock.go.


func TestService(t *testing.T) {
	// creates the map using the make function and assigns it to the users field in 
	// the instance of mockRepository
	mockRepo := &mockRepository{make(map[uuid.UUID]entity.User)}

	// Create a new instance of the service using a mockRepo instance
	service := NewService(mockRepo)
	assert.NotNil(t, service)

	// A new instance of CreateUserInput struct
	input := CreateUserInput{
		Name:     "Name",
		Email:    "testemail@gmail.com",
		Password: "password123!",
		IsGuest:  false,
	}
	t.Run("Create user", func(t *testing.T) {
		t.Run("with a valid user", func(t *testing.T) {

			// Call the CreateUser function in service.go on our service instance 
			user, err := service.CreateUser(input)
			assert.NoError(t, err)
			assert.Equal(t, input.Name, user.Name)
		})

		t.Run("with an invalid email", func(t *testing.T) {
			invalidInput := input
			invalidInput.Email = "xyz.com"
			_, err := service.CreateUser(invalidInput)
			assert.ErrorIs(t, err, ErrInvalidEmail)
		})

	})
}