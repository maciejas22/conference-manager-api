package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/maciejas22/conference-manager-api/cm-info/config"
	"github.com/maciejas22/conference-manager-api/cm-info/internal/db"
	handler "github.com/maciejas22/conference-manager-api/cm-info/internal/handler"
	newsRepo "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/news"
	tosRepo "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/tos"
	newsService "github.com/maciejas22/conference-manager-api/cm-info/internal/service/news"
	tosService "github.com/maciejas22/conference-manager-api/cm-info/internal/service/tos"
	"google.golang.org/grpc"
)

func initNewsService(dbClient *db.DB) newsService.NewsServiceInterface {
	newsRepo := newsRepo.NewNewsRepo(dbClient.Conn)
	return newsService.NewNewsService(newsRepo)
}

func initToSService(dbClient *db.DB) tosService.ToSServiceInterface {
	tosRepo := tosRepo.NewToSRepo(dbClient.Conn)
	return tosService.NewToSService(tosRepo)
}

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

	newsService := initNewsService(db)
	tosService := initToSService(db)

	grpcServer := grpc.NewServer()
	handler.NewServer(grpcServer, newsService, tosService)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("APP", "failed to serve gRPC server over port", config.AppConfig.Port, "error", err)
	}
}
