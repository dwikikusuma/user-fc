package usecase

import (
	"commerce/cmd/user/service"
	"commerce/infrastructure/log"
	"commerce/models"
	"commerce/utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserUseCase struct {
	UserService service.UserService
	JwtSecrete  string
}

func NewUserUseCase(userService service.UserService, jwtSecrete string) *UserUseCase {
	return &UserUseCase{
		UserService: userService,
		JwtSecrete:  jwtSecrete,
	}

}
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := uc.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, user models.User) error {
	hashedPass, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": user.Email,
		}).Errorf("utils.GenerateHashPassword() got error %v", err)
		return err
	}
	user.Password = hashedPass
	_, err = uc.UserService.CreateNewUser(ctx, user)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": user.Email,
			"name":  user.Name,
		}).Errorf("uc.UserService.CreateNewUser(user) got error %v", err)
		return err
	}

	return nil
}

func (uc *UserUseCase) Login(ctx context.Context, loginRequest models.LoginParameter) (string, error) {
	user, err := uc.UserService.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": loginRequest.Email,
		}).Errorf("uc.UserService.GetUserByEmail %v", err)
	}

	isMatch, _ := utils.ComparePassword(loginRequest.Password, user.Password)

	if !isMatch {
		return "", errors.New("invalid username or password")
	}

	currentTime := time.Now()
	expTime := currentTime.Add(time.Hour * 24)

	// Log the current time and expiration time
	log.Logger.WithFields(logrus.Fields{
		"user_id":      user.ID,
		"current_time": currentTime,
		"exp_time":     expTime,
		"exp_unix":     expTime.Unix(), // Log the Unix timestamp of the expiration
	}).Info("Token created")

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expTime.Unix(),
	})

	token, err := tokenString.SignedString([]byte(uc.JwtSecrete))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"email": loginRequest.Email,
		}).Errorf("token.SignedString() got an error %v", err)
		return "", err
	}

	return token, nil
}

func (uc *UserUseCase) GetUserById(ctx context.Context, userId int64) (*models.User, error) {
	user, err := uc.UserService.GetUserByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
