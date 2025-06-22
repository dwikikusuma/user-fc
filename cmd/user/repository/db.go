package repository

import (
	"commerce/models"
	"context"
	"fmt"
)

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.Database.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) InsertUser(ctx context.Context, user *models.User) (int64, error) {
	err := r.Database.WithContext(ctx).Create(user).Error
	fmt.Println("InsertUser error:", err)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepository) FindByUserId(ctx context.Context, userId int64) (*models.User, error) {
	var user models.User
	err := r.Database.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
