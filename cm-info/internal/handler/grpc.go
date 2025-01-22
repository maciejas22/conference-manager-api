package grpc

import (
	"context"

	f "github.com/maciejas22/conference-manager-api/cm-info/internal/service/model"
	news "github.com/maciejas22/conference-manager-api/cm-info/internal/service/news"
	tos "github.com/maciejas22/conference-manager-api/cm-info/internal/service/tos"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/info"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(grpcServer *grpc.Server, newsService news.NewsServiceInterface, tosService tos.ToSServiceInterface) {
	s := &InfoServ{
		newsService: newsService,
		tosService:  tosService,
	}
	pb.RegisterInfoServiceServer(grpcServer, s)
}

func (srv *InfoServ) ListNews(ctx context.Context, req *pb.NewsRequest) (*pb.NewsResponse, error) {
	page := f.Page{
		PageSize:   int(req.Page.Size),
		PageNumber: int(req.Page.Number),
	}
	news, meta, err := srv.newsService.GetNews(page)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	n := make([]*pb.NewsItem, len(news))
	for i, item := range news {
		n[i] = &pb.NewsItem{
			Id:      int32(item.Id),
			Title:   item.Title,
			Content: item.Content,
			Date:    timestamppb.New(item.CreatedAt),
		}
	}

	return &pb.NewsResponse{
		NewsPage: &pb.NewsPage{
			Data: n,
			Meta: &pb.PageInfo{
				TotalItems: int32(meta.TotalItems),
				TotalPages: int32(meta.TotalPages),
				Size:       int32(meta.PageSize),
				Number:     int32(meta.PageNumber),
			},
		},
	}, nil
}

func (srv *InfoServ) GetTermsOfService(ctx context.Context, req *emptypb.Empty) (*pb.TermsOfServiceResponse, error) {
	tos, err := srv.tosService.GetTermsOfService()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tgtSections := make([]*pb.Section, len(tos.Sections))
	for i, sec := range tos.Sections {
		if sec != nil {
			tgtSubsections := make([]*pb.SubSection, len(sec.Subsections))
			for j, sub := range sec.Subsections {
				if sub != nil {
					tgtSubsections[j] = &pb.SubSection{
						Id:      int32(sub.Id),
						Title:   sub.Title,
						Content: *sub.Content,
					}
				}
			}

			tgtSections[i] = &pb.Section{
				Id:          int32(sec.Id),
				Title:       sec.Title,
				Content:     *sec.Content,
				Subsections: tgtSubsections,
			}
		}
	}

	return &pb.TermsOfServiceResponse{
		TermsOfService: &pb.TermsOfService{
			Id:              int32(tos.Id),
			Introduction:    tos.Introduction,
			Acknowledgement: tos.Acknowledgement,
			Sections:        tgtSections,
		},
	}, nil
}
