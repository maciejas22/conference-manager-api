package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/maciejas22/conference-manager-api/cm-conferences/config"
	"github.com/maciejas22/conference-manager-api/cm-conferences/internal/db"
	handler "github.com/maciejas22/conference-manager-api/cm-conferences/internal/handler"
	"github.com/maciejas22/conference-manager-api/cm-conferences/internal/middleware"

	ra "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/agenda"
	rc "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference"
	rco "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference_organizer"
	rcp "github.com/maciejas22/conference-manager-api/cm-conferences/internal/repository/conference_participant"

	sa "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/agenda"
	sc "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/conference"
	sco "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/organizer"
	scp "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/participant"
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

	conferenceRepo := rc.NewConferenceRepo(db.Conn)
	organizerRepo := rco.NewOrganizerRepo(db.Conn)
	participantRepo := rcp.NewParticipantRepo(db.Conn)
	agendaRepo := ra.NewAgendaRepo(db.Conn)

	conferenceService := sc.NewConferenceService(conferenceRepo, agendaRepo)
	organizerService := sco.NewConferenceOrganizerService(organizerRepo)
	participantService := scp.NewParticipantService(participantRepo)
	agendaService := sa.NewAgendaService(agendaRepo)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryErrorInterceptor(logger)),
	)
	handler.NewServer(grpcServer, conferenceService, organizerService, participantService, agendaService)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("APP", "failed to serve gRPC server over port", config.AppConfig.Port, "error", err)
	}
}
