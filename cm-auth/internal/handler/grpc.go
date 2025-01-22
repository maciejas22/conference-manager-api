package grpc

import (
	"context"

	service "github.com/maciejas22/conference-manager-api/cm-auth/internal/service/auth"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer(grpcServer *grpc.Server, service service.AuthServiceInterface) {
	s := &AuthServ{service: service}
	pb.RegisterAuthServiceServer(grpcServer, s)
}

func (srv *AuthServ) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var stripeAccountId *string
	if req.StripeAccountId != "" {
		stripeAccountId = &req.StripeAccountId
	}

	var role service.Role
	if req.Role == pb.Role_ROLE_ORGANIZER {
		role = service.Organizer
	} else {
		role = service.Participant
	}

	id, err := srv.service.RegisterUser(service.RegisterUserInput{
		Email:           req.Email,
		Password:        req.Password,
		Role:            role,
		StripeAccountId: stripeAccountId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		UserId: int32(*id),
	}, nil
}

func (srv *AuthServ) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	session, err := srv.service.LoginUser(service.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		SessionId: *session,
	}, nil
}

func (srv *AuthServ) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	_, err := srv.service.LogoutUser(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LogoutResponse{}, nil
}

func (srv *AuthServ) UserProfileBySession(ctx context.Context, req *pb.UserProfileBySessionRequest) (*pb.UserProfileBySessionResponse, error) {
	user, err := srv.service.GetUserBySession(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var role pb.Role
	if user.Role == service.Organizer {
		role = pb.Role_ROLE_ORGANIZER
	} else {
		role = pb.Role_ROLE_PARTICIPANT
	}
	var userName, userSurname, userUsername, userStripeAccountId string
	if user.Name != nil {
		userName = *user.Name
	}
	if user.Surname != nil {
		userSurname = *user.Surname
	}
	if user.Username != nil {
		userUsername = *user.Username
	}
	if user.StripeAccountId != nil {
		userStripeAccountId = *user.StripeAccountId
	}
	return &pb.UserProfileBySessionResponse{
		User: &pb.User{
			UserId:          int32(user.Id),
			Email:           user.Email,
			Role:            role,
			Name:            userName,
			Surname:         userSurname,
			Username:        userUsername,
			StripeAccountId: userStripeAccountId,
		},
	}, nil
}

func (srv *AuthServ) UpdateSession(ctx context.Context, req *pb.UpdateSessionRequest) (*pb.UpdateSessionResponse, error) {
	s, uId, err := srv.service.UpdateSession(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateSessionResponse{
		SessionId: s,
		UserId:    int32(uId),
	}, nil
}

func (srv *AuthServ) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	var userInput service.UpdateUserInput
	if req.Email != "" {
		userInput.Email = &req.Email
	}
	if req.Name != "" {
		userInput.Name = &req.Name
	}
	if req.Surname != "" {
		userInput.Surname = &req.Surname
	}
	if req.Username != "" {
		userInput.Username = &req.Username
	}
	if req.Password != "" {
		userInput.Password = &req.Password
	}
	_, err := srv.service.UpdateUser(int(req.UserId), userInput)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateUserProfileResponse{
		UserId: int32(req.UserId),
	}, nil
}

func (srv *AuthServ) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	ok, err := srv.service.VerifySession(req.SessionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ValidateSessionResponse{
		IsValid: ok,
	}, nil
}

func (srv *AuthServ) UserProfileById(ctx context.Context, req *pb.UserProfileByIdRequest) (*pb.UserProfileByIdResponse, error) {
	user, err := srv.service.GetUser(int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var role pb.Role
	if user.Role == service.Organizer {
		role = pb.Role_ROLE_ORGANIZER
	} else {
		role = pb.Role_ROLE_PARTICIPANT
	}
	var userName, userSurname, userUsername, userStripeAccountId string
	if user.Name != nil {
		userName = *user.Name
	}
	if user.Surname != nil {
		userSurname = *user.Surname
	}
	if user.Username != nil {
		userUsername = *user.Username
	}
	if user.StripeAccountId != nil {
		userStripeAccountId = *user.StripeAccountId
	}
	return &pb.UserProfileByIdResponse{
		User: &pb.User{
			UserId:          int32(user.Id),
			Email:           user.Email,
			Role:            role,
			Name:            userName,
			Surname:         userSurname,
			Username:        userUsername,
			StripeAccountId: userStripeAccountId,
		},
	}, nil
}
