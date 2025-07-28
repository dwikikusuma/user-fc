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
	"google.golang.org/grpc/reflection"
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

	log.Logger.Info("Setting up routes...")
	routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)

	//_ = router.Run(":" + port)
	log.Logger.Info("Server started on port " + port)

	grpcPort := ":5051"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Logger.Fatalf("Failed to listen on %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &userGrpc.GRPCServer{UserUseCase: *userUseCase})
	reflection.Register(grpcServer)

	log.Logger.Infof("gRPC server started on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Logger.Fatalf("Failed to start gRPC server: %v", err)
	}
}
