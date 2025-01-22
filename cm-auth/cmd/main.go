package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/maciejas22/conference-manager-api/cm-auth/config"
	"github.com/maciejas22/conference-manager-api/cm-auth/internal/db"
	handler "github.com/maciejas22/conference-manager-api/cm-auth/internal/handler"
	"github.com/maciejas22/conference-manager-api/cm-auth/internal/middleware"
	sessionRepository "github.com/maciejas22/conference-manager-api/cm-auth/internal/repository/session"
	userRepository "github.com/maciejas22/conference-manager-api/cm-auth/internal/repository/user"
	service "github.com/maciejas22/conference-manager-api/cm-auth/internal/service/auth"
	"google.golang.org/grpc"
)

func initLogger(config *config.Config) *slog.Logger {
	var level slog.Level
	if config.GoEnv == "dev" {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	return logger
}

func initAuthService(dbClient *db.DB) service.AuthServiceInterface {
	userRepo := userRepository.NewUserRepo(dbClient.Conn)
	sessionRepo := sessionRepository.NewSessionRepo(dbClient.Conn)

	return service.NewAuthService(userRepo, sessionRepo)
}

func main() {
	ctx := context.Background()
	config.Init()

	logger := initLogger(config.AppConfig)

	db, err := db.Connect(ctx, logger)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AppConfig.Port))
	if err != nil {
		logger.Error("APP", "failed to listen on port", config.AppConfig.Port, "error", err)
	}
	logger.Info("APP", "listening on port", config.AppConfig.Port)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryErrorInterceptor(logger)),
	)

	authService := initAuthService(db)
	handler.NewServer(grpcServer, authService)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("APP", "failed to serve gRPC server over port", config.AppConfig.Port, "error", err)
	}
}
