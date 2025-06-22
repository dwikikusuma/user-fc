package service

import (
	"commerce/cmd/user/repository"
	"commerce/models"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (svc *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := svc.UserRepo.FindByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *UserService) CreateNewUser(ctx context.Context, user models.User) (int64, error) {
	userId, err := svc.UserRepo.InsertUser(ctx, &user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (svc *UserService) GetUserByUserId(ctx context.Context, userId int64) (*models.User, error) {
	user, err := svc.UserRepo.FindByUserId(ctx, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
