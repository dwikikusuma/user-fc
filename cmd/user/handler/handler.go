package handler

import (
	"commerce/cmd/user/usecase"
	"commerce/infrastructure/log"
	"commerce/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	// User Dependency
	UserUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}
}

func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var parameter models.RegisterParameter
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&parameter); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if len(parameter.Password) < 6 {
		c.JSON(400, gin.H{
			"error": "Password must be at least 6 characters long",
		})
		return
	}

	if parameter.Password != parameter.ConfirmPassword {
		c.JSON(400, gin.H{
			"error": "Passwords do not match",
		})
		return
	}

	user, err := h.UserUseCase.GetUserByEmail(ctx, parameter.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if user != nil && user.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "email already exists",
		})
		return
	}

	err = h.UserUseCase.RegisterUser(ctx, models.User{
		Name:     parameter.Name,
		Email:    parameter.Email,
		Password: parameter.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginParam models.LoginParameter
	ctx := c.Request.Context()

	err := c.ShouldBindJSON(&loginParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input request",
		})
		return
	}

	token, err := h.UserUseCase.Login(ctx, loginParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to login",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userIDStr, isExist := c.Get("user_id")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_message": "Unauthorized",
		})
		return
	}

	userID, ok := userIDStr.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_message": "Invalid user id",
		})
		return
	}

	ctx := c.Request.Context()

	user, err := h.UserUseCase.GetUserById(ctx, int64(userID))
	if err != nil || user == nil {
		log.Logger.WithFields(logrus.Fields{
			"user_id": userID,
		}).Errorf("h.UserUseCase.GetUserById() got an error %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to get user info",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"email": user.Email,
	})
}
