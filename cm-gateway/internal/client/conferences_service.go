package client

import (
	"context"
	"log"

	"github.com/maciejas22/conference-manager-api/cm-gateway/config"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/conference"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitConferencesClient(ctx context.Context) pb.ConferenceServiceClient {
	conn, err := grpc.NewClient(config.AppConfig.ConferenceServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create conferences service client: %v", err)
	}

	return pb.NewConferenceServiceClient(conn)
}
