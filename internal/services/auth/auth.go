package services

import (
	"context"
	"errors"
	"fmt"
	"lesson-proj/internal/database"
	"lesson-proj/internal/models"
	authUtils "lesson-proj/internal/services/auth/utils"
	"os"
)

type UserService struct {
	repository *database.UserRepository
}

func NewUserService(repository *database.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (service *UserService) Registration(ctx context.Context, input models.CreateUser) (*models.UserWithoutPassword, error) {
	if err := authUtils.ValidateCreateUserInput(input.Email, input.Name, input.Password); err != nil {
		return nil, err
	}

	
	hashPassword, err := authUtils.HashPassword(input.Password, os.Getenv("PASSWORD_PEPPER"))
	if err != nil {
		return nil, err
	}
	// avoid keeping plaintext longer than needed
	input.Password = ""

	createdUser, err := service.repository.CreateUser(ctx, models.CreateUser{
		Email:    input.Email,
		Name:     input.Name,
		Password: hashPassword,
	})
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}


func (service *UserService) Authorization(ctx context.Context, email, password string) (*models.UserWithoutPassword, error) {
	user, err := service.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}
	pepper := os.Getenv("PASSWORD_PEPPER")
	if pepper == "" {
		return nil, errors.New("password pepper is not configured")
	}
	ok, err := authUtils.VerifyPassword(password, pepper, user.HashedPassword)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("invalid password")
	}
	userWithoutPassword := &models.UserWithoutPassword{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	return userWithoutPassword, nil
	
}

func (service *UserService) GetAllUsers(ctx context.Context) ([]models.UserWithoutPassword, error) {
	users, err := service.repository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service *UserService) GetUserByID(ctx context.Context, id int) (*models.UserWithoutPassword, error) {
	user, err := service.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) UpdateUser(ctx context.Context, id int, input models.UpdateUser) (*models.UserWithoutPassword, error) {
	if err := authUtils.ValidateUpdateUserInput(input.Email, input.Name, input.Password); err != nil {
		return nil, err
	}
	updatedUser, err := service.repository.UpdateUser(ctx, id, input)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (service *UserService) DeleteUser(ctx context.Context, id int) error {
	return service.repository.DeleteUser(ctx, id)
}
