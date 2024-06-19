package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/leegodden/boards/server/db"
	"github.com/leegodden/boards/server/internal/entity"
)

type Repository interface {
	CreateUser(user entity.User) error
	DeleteUser(userId uuid.UUID) error
}

type repository struct {
	db *db.DB
}

func (r *repository) CreateUser(user entity.User) error {
	ctx := context.Background()
	sql := "INSERT INTO users (id, name, email, password, is_guest, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, sql, user.Id, user.Name, user.Email, user.Password, user.IsGuest, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteUser(userId uuid.UUID) error {
	ctx := context.Background()
	sql := "DELETE from users where id = $1"
	_, err := r.db.Exec(ctx, sql, userId)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *db.DB) Repository {
	return &repository{db}
}

