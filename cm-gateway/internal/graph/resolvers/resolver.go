package resolvers

import (
	"context"
	"log/slog"

	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/client"
	authPb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
	conferencePb "github.com/maciejas22/conference-manager-api/cm-proto/conference"
	infoPb "github.com/maciejas22/conference-manager-api/cm-proto/info"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ctx                     context.Context
	authServiceClient       authPb.AuthServiceClient
	infoServiceClient       infoPb.InfoServiceClient
	conferenceServiceClient conferencePb.ConferenceServiceClient
	s3Client                client.S3Service
}

func NewResolver(ctx context.Context, logger *slog.Logger, authServiceClient authPb.AuthServiceClient, conferenceServiceClient conferencePb.ConferenceServiceClient) *Resolver {
	return &Resolver{
		ctx:                     ctx,
		authServiceClient:       authServiceClient,
		infoServiceClient:       client.InitInfoClient(ctx),
		conferenceServiceClient: conferenceServiceClient,
		s3Client:                client.NewS3Client(ctx, logger),
	}
}
