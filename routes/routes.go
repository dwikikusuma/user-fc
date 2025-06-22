package routes

import (
	"commerce/cmd/user/handler"
	"commerce/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler handler.UserHandler, jwtSecrete string) {
	// Public
	router.Use(middleware.RequestLogger())
	router.GET("/ping", userHandler.Ping)
	router.POST("/v1/register", userHandler.Register)
	router.POST("/v1/login", userHandler.Login)

	//Private
	authMiddleware := middleware.AuthMiddleWare(jwtSecrete)
	private := router.Group("/api")
	private.Use(authMiddleware)
	private.GET("/v1/user_info", userHandler.GetUserInfo)
}
