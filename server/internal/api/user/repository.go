package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/leegodden/boards/server/db"
	"github.com/leegodden/boards/server/internal/entity"
)

var (
	ErrDatabase              = errors.New("database error")
	ErrUniqueEmailConstraint = errors.New("not a unique email")
)


type Repository interface {
	CreateUser(user entity.User) error
	DeleteUser(userId uuid.UUID) error
}

type repository struct {
	db *db.DB
}

// Associate repository struct with Repository interface accepts a user entity 
// and inserts it to the database. Also specifically handles unique email 
// constraint violation.
func (r *repository) CreateUser(user entity.User) error {
	ctx := context.Background()
	sql := "INSERT INTO users (id, name, email, password, is_guest, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, sql, user.Id, user.Name, user.Email, user.Password, user.IsGuest, user.CreatedAt, user.UpdatedAt)


// If any error occurs, the error is further examined to check if it's a 
// PgError (a specific kind of error from the PostgreSQL database).
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation && e.ConstraintName == "users_email_key" {
			return ErrUniqueEmailConstraint
		}
		return fmt.Errorf("%v: %w", ErrDatabase, err)
	}
	return nil
}

func (r *repository) DeleteUser(userId uuid.UUID) error {
	ctx := context.Background()
	sql := "DELETE from users where id = $1"
	_, err := r.db.Exec(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrDatabase, err)
	}
	return nil
}

func NewRepository(db *db.DB) Repository {
	return &repository{db}
}

