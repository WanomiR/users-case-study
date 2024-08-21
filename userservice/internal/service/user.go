package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wanomir/e"
	"userservice/internal/entity"
	"userservice/internal/infrastructure/repository"
)

type UserService struct {
	DB repository.Repository
}

func NewUserService(db repository.Repository) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) GetUsers(ctx context.Context, offset, limit int) (users []entity.User, err error) {
	if limit < 1 {
		return []entity.User{}, errors.New("limit must be greater than zero")
	}

	users, err = s.DB.ListUsers(ctx, offset, limit)
	if err != nil {
		return []entity.User{}, e.WrapIfErr("failed to get users", err)
	}

	return users, err
}

func (s *UserService) GetUser(ctx context.Context, email string) (entity.User, error) {
	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.User{}, e.WrapIfErr("failed to get user by email", err)
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user entity.User) (int, error) {
	if !user.AgeIsValid() {
		return 0, errors.New(fmt.Sprintf("user age must be greater than 17, got %v", user.Age))
	}

	if _, err := s.DB.GetUserByEmail(ctx, user.Email); err == nil {
		return 0, e.WrapIfErr("user with this email already exists", err)
	}

	userId, err := s.DB.CreateUser(ctx, user)
	if err != nil {
		return 0, e.WrapIfErr("failed to create user", err)
	}

	return userId, nil
}
