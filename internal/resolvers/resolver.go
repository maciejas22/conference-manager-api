package resolvers

import (
	"context"

	"github.com/maciejas22/conference-manager/api/internal/db"
	"github.com/maciejas22/conference-manager/api/pkg/s3"
)

type Resolver struct {
	ctx      context.Context
	dbClient *db.DB
	s3Client *s3.S3Client
}

func NewResolver(ctx context.Context, db *db.DB, s3 *s3.S3Client) *Resolver {
	return &Resolver{
		ctx:      ctx,
		dbClient: db,
		s3Client: s3,
	}
}
