package client

import (
	"context"
	"log"

	"github.com/maciejas22/conference-manager-api/cm-gateway/config"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthClient(ctx context.Context) pb.AuthServiceClient {
	conn, err := grpc.NewClient(config.AppConfig.AuthServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create auth service client: %v", err)
	}

	return pb.NewAuthServiceClient(conn)
}
