package repository

import (
	"context"
	"database/sql"
	"userservice/internal/entity"
)

type Repository interface {
	Connection() *sql.DB
	ListUsers(ctx context.Context, offset, limit int) ([]entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (int, error)
	UpdateUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, email string) error
}
