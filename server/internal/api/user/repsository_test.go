package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leegodden/boards/server/internal/entity"
	"github.com/leegodden/boards/server/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)
	repo := NewRepository(db)

	testUser := newTestUser()

	t.Run("Create user", func(t *testing.T) {
		err := repo.CreateUser(testUser)
		assert.NoError(t, err)
	})

	t.Run("Delete user", func(t *testing.T) {
		err := repo.DeleteUser(testUser.Id)
		assert.NoError(t, err)
	})

	t.Run("Create user with non-unique email", func(t *testing.T) {
		testUserBadEmail := testUser
		testUserBadEmail.Id = uuid.New()
		err := repo.CreateUser(testUserBadEmail)
		assert.ErrorIs(t, err, ErrUniqueEmailConstraint)
	})
}

func newTestUser() entity.User {
	email := uuid.New().String() + "email.com"
	user := entity.User{
		Id:        uuid.New(),
		Name:      "testname",
		Email:     email,
		Password:  "password123!",
		IsGuest:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user
}