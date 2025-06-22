package main

import (
	"commerce/cmd/user/handler"
	"commerce/cmd/user/repository"
	"commerce/cmd/user/resource"
	"commerce/cmd/user/service"
	"commerce/cmd/user/usecase"
	"commerce/config"
	"commerce/infrastructure/log"
	"commerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redisClient := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)
	log.SetupLogger()
	port := cfg.App.Port
	router := gin.Default()

	userRepository := repository.NewUserRepository(redisClient, db)
	userService := service.NewUserService(*userRepository)
	userUseCase := usecase.NewUserUseCase(*userService, cfg.Secret.JWTSecret)
	userHandler := handler.NewUserHandler(*userUseCase)
	routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)
	_ = router.Run(":" + port)
	log.Logger.Info("Server started on port " + port)
}
