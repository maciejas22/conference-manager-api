package grpc

import (
	service "github.com/maciejas22/conference-manager-api/cm-auth/internal/service/auth"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
)

type AuthServ struct {
	service service.AuthServiceInterface
	pb.UnimplementedAuthServiceServer
}
