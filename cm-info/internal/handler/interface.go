package grpc

import (
	ns "github.com/maciejas22/conference-manager-api/cm-info/internal/service/news"
	ts "github.com/maciejas22/conference-manager-api/cm-info/internal/service/tos"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/info"
)

type InfoServ struct {
	newsService ns.NewsServiceInterface
	tosService  ts.ToSServiceInterface
	pb.UnimplementedInfoServiceServer
}
