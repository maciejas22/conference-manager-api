package client

import (
	"context"
	"log"

	"github.com/maciejas22/conference-manager-api/cm-gateway/config"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/info"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInfoClient(ctx context.Context) pb.InfoServiceClient {
	conn, err := grpc.NewClient(config.AppConfig.InfoServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create info service client: %v", err)
	}

	return pb.NewInfoServiceClient(conn)
}
