package user

import (
	"github.com/google/uuid"
	"github.com/leegodden/boards/server/internal/entity"
)

// A mock repository that uses a map to store the users
// with UUID as the key and the "User" struct as the values

type mockRepository struct {
	users map[uuid.UUID]entity.User
}

func (r *mockRepository) CreateUser(user entity.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *mockRepository) DeleteUser(userId uuid.UUID) error {
	delete(r.users, userId)
	return nil
}

