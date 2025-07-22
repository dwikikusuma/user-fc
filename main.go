package main

import (
	"commerce/cmd/user/handler"
	"commerce/cmd/user/repository"
	"commerce/cmd/user/resource"
	"commerce/cmd/user/service"
	"commerce/cmd/user/usecase"
	"commerce/config"
	userGrpc "commerce/grpc"
	"commerce/infrastructure/log"
	"commerce/proto/userpb"
	"commerce/routes"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
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

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &userGrpc.GRPCServer{UserUseCase: *userUseCase})

	lis, _ := net.Listen("tcp", ":50051")
	err := grpcServer.Serve(lis)
	if err != nil {
		return
	}

	routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)
	_ = router.Run(":" + port)

	log.Logger.Info("Server started on port " + port)
	log.Logger.Info("gRPC server started on port 50051")
}
