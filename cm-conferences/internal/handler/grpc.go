package grpc

import (
	"context"
	"log"
	"time"

	a "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/agenda"
	commonService "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/common"
	c "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/conference"
	cp "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/organizer"
	co "github.com/maciejas22/conference-manager-api/cm-conferences/internal/service/participant"
	commonPb "github.com/maciejas22/conference-manager-api/cm-proto/common"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/conference"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewServer(grpcServer *grpc.Server,
	conferenceService c.ConferenceServiceInterface,
	organizerService cp.OrganizerServiceInterface,
	participantService co.ParticipantServiceInterface,
	agendaService a.AgendaServiceInterface,
) {
	s := &ConferenceServ{
		conferenceService:  conferenceService,
		organizerService:   organizerService,
		participantService: participantService,
		agendaService:      agendaService,
	}
	pb.RegisterConferenceServiceServer(grpcServer, s)
}

func (srv *ConferenceServ) GetConferences(ctx context.Context, req *pb.ConferencesRequest) (*pb.ConferencesResponse, error) {
	cIds := make([]int, len(req.ConferenceIds))
	for i, id := range req.ConferenceIds {
		cIds[i] = int(id)
	}
	c, err := srv.conferenceService.GetConferencesByIds(ctx, cIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	aCount, err := srv.agendaService.GetAgendaItemsCount(cIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	aCountMap := make(map[int]int)
	for _, count := range aCount {
		aCountMap[count.ConferenceId] = count.AgendaItemsCount
	}

	pCount, err := srv.participantService.GetParticipantsCount(cIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pCountMap := make(map[int]int)
	for _, count := range pCount {
		pCountMap[count.ConferenceId] = count.ParticipantsCount
	}

	conferences := make([]*pb.Conference, len(c))
	for i, conf := range c {
		var participantsLimit *int32
		if conf.ParticipantsLimit != nil {
			limit := int32(*conf.ParticipantsLimit)
			participantsLimit = &limit
		}

		var registrationDeadline *timestamppb.Timestamp
		if conf.RegistrationDeadline != nil {
			registrationDeadline = timestamppb.New(*conf.RegistrationDeadline)
		}

		conferences[i] = &pb.Conference{
			Id:                   int32(conf.Id),
			Title:                conf.Title,
			StartDate:            timestamppb.New(conf.StartDate),
			EndDate:              timestamppb.New(conf.EndDate),
			Location:             conf.Location,
			Website:              *conf.Website,
			Acronym:              *conf.Acronym,
			AdditionalInfo:       *conf.AdditionalInfo,
			RegistrationDeadline: registrationDeadline,
			TicketPrice:          int32(conf.TicketPrice),
			ParticipantsLimit:    *participantsLimit,
			ParticipantsCount:    int32(pCountMap[conf.Id]),
			EventsCount:          int32(aCountMap[conf.Id]),
		}
	}

	return &pb.ConferencesResponse{
		Conferences: conferences,
	}, nil
}
func (srv *ConferenceServ) GetConferenceAgenda(ctx context.Context, req *pb.AgendaRequest) (*pb.AgendaResponse, error) {
	a, err := srv.agendaService.GetAgenda(int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	agendaItems := make([]*pb.AgendaItem, len(a))
	for i, item := range a {
		agendaItems[i] = &pb.AgendaItem{
			Id:        int32(item.Id),
			Event:     item.Event,
			Speaker:   item.Speaker,
			StartTime: timestamppb.New(item.StartTime),
			EndTime:   timestamppb.New(item.EndTime),
		}
	}

	return &pb.AgendaResponse{
		Agenda: agendaItems,
	}, nil
}

func (srv *ConferenceServ) GetConferencesPage(ctx context.Context, req *pb.ConferencesPageRequest) (*pb.ConferencesPageResponse, error) {
	var order commonService.Order
	if req.Sort.Order == commonPb.Order_ASC {
		order = commonService.ASC
	} else {
		order = commonService.DESC
	}
	c, m, err := srv.conferenceService.GetConferencesPage(ctx, int(req.UserId), &commonService.Page{
		PageNumber: int(req.Page.Number),
		PageSize:   int(req.Page.Size),
	}, &commonService.Sort{
		Column: req.Sort.Column,
		Order:  order,
	}, &c.ConferencesFilters{
		Title:          &req.Filters.Title,
		AssociatedOnly: &req.Filters.AssociatedOnly,
		RunningOnly:    &req.Filters.RunningOnly,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cIds []int32
	for _, id := range c {
		cIds = append(cIds, int32(id))
	}

	return &pb.ConferencesPageResponse{
		ConferencesPage: &pb.ConferencesPage{
			ConferenceIds: cIds,
			Meta: &pb.ConferenceMeta{
				Page: &commonPb.PageInfo{
					Number:     int32(m.PageNumber),
					Size:       int32(m.PageSize),
					TotalItems: int32(m.TotalItems),
					TotalPages: int32(m.TotalPages),
				},
			},
		},
	}, nil
}

func (srv *ConferenceServ) GetConferencesMetrics(ctx context.Context, req *pb.ConferencesMetricsRequest) (*pb.ConferencesMetricsResponse, error) {
	m, err := srv.conferenceService.GetConferencesMetrics(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConferencesMetricsResponse{
		Metrics: &pb.ConferencesMetrics{
			RunningConferences:        int32(m.RunningConferences),
			StartingInLessThan24Hours: int32(m.StartingInLessThan24Hours),
			TotalConducted:            int32(m.TotalConducted),
			ParticipantsToday:         int32(m.ParticipantsToday),
		},
	}, nil
}
func (srv *ConferenceServ) CreateConference(ctx context.Context, req *pb.CreateConferenceRequest) (*pb.CreateConferenceResponse, error) {
	var participantsLimit *int
	if req.CreateConferenceInput.ParticipantsLimit != 0 {
		limit := int(req.CreateConferenceInput.ParticipantsLimit)
		participantsLimit = &limit
	} else {
		participantsLimit = nil
	}
	var registrationDeadline *time.Time
	if req.CreateConferenceInput.RegistrationDeadline != nil {
		rd := req.CreateConferenceInput.RegistrationDeadline.AsTime()
		registrationDeadline = &rd
	}
	a := make([]*c.AgendaItem, len(req.CreateConferenceInput.Agenda))
	for i, item := range req.CreateConferenceInput.Agenda {
		a[i] = &c.AgendaItem{
			Id:        0,
			Event:     item.Event,
			Speaker:   item.Speaker,
			StartTime: item.StartTime.AsTime(),
			EndTime:   item.EndTime.AsTime(),
		}
	}

	c, err := srv.conferenceService.CreateConference(ctx, int(req.UserId), c.CreateConferenceInput{
		Title:                req.CreateConferenceInput.Title,
		StartDate:            req.CreateConferenceInput.StartDate.AsTime(),
		EndDate:              req.CreateConferenceInput.EndDate.AsTime(),
		Location:             req.CreateConferenceInput.Location,
		Website:              &req.CreateConferenceInput.Website,
		Acronym:              &req.CreateConferenceInput.Acronym,
		AdditionalInfo:       &req.CreateConferenceInput.AdditionalInfo,
		ParticipantsLimit:    participantsLimit,
		RegistrationDeadline: registrationDeadline,
		TicketPrice:          int(req.CreateConferenceInput.TicketPrice),
		Agenda:               a,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = srv.organizerService.AddOrganizerToConference(int(req.UserId), *c)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateConferenceResponse{
		ConferenceId: int32(*c),
	}, nil
}
func (srv *ConferenceServ) ModifyConference(ctx context.Context, req *pb.ModifyConferenceRequest) (*pb.ModifyConferenceResponse, error) {
	startDate := req.Input.StartDate.AsTime()
	endDate := req.Input.EndDate.AsTime()
	var participantsLimit *int
	if req.Input.ParticipantsLimit != 0 {
		limit := int(req.Input.ParticipantsLimit)
		participantsLimit = &limit
	} else {
		participantsLimit = nil
	}
	var registrationDeadline *time.Time
	if req.Input.RegistrationDeadline != nil {
		rd := req.Input.RegistrationDeadline.AsTime()
		registrationDeadline = &rd
	}
	ticketPrice := int(req.Input.TicketPrice)

	a := make([]*c.AgendaItem, len(req.Input.Agenda))
	for i, item := range req.Input.Agenda {
		if item.DeleteItem != 0 {
			a[i] = &c.AgendaItem{
				Id:        int(item.DeleteItem),
				Event:     "",
				Speaker:   "",
				StartTime: time.Now(),
				EndTime:   time.Now(),
			}
		} else {
			a[i] = &c.AgendaItem{
				Id:        0,
				Event:     item.CreateItem.Event,
				Speaker:   item.CreateItem.Speaker,
				StartTime: item.CreateItem.StartTime.AsTime(),
				EndTime:   item.CreateItem.EndTime.AsTime(),
			}
		}
	}

	cId, err := srv.conferenceService.ModifyConference(ctx, c.ModifyConferenceInput{
		ID:                   int(req.Input.Id),
		Title:                &req.Input.Title,
		StartDate:            &startDate,
		EndDate:              &endDate,
		Location:             &req.Input.Location,
		Website:              &req.Input.Website,
		Acronym:              &req.Input.Acronym,
		AdditionalInfo:       &req.Input.AdditionalInfo,
		ParticipantsLimit:    participantsLimit,
		RegistrationDeadline: registrationDeadline,
		TicketPrice:          &ticketPrice,
		Agenda:               a,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ModifyConferenceResponse{
		ConferenceId: int32(*cId),
	}, nil
}
func (srv *ConferenceServ) AddUserToConference(ctx context.Context, req *pb.AddUserToConferenceRequest) (*pb.AddUserToConferenceResponse, error) {
	tId, err := srv.participantService.AddUserToConference(int(req.UserId), int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.AddUserToConferenceResponse{
		TicketId: tId,
	}, nil
}

func (srv *ConferenceServ) RemoveUserFromConference(ctx context.Context, req *pb.RemoveUserFromConferenceRequest) (*pb.RemoveUserFromConferenceResponse, error) {
	_, err := srv.participantService.RemoveUserFromConference(int(req.UserId), int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveUserFromConferenceResponse{}, nil
}
func (srv *ConferenceServ) IsUserAssociatedWithConference(ctx context.Context, req *pb.IsUserAssociatedWithConferenceRequest) (*pb.IsUserAssociatedWithConferenceResponse, error) {
	isParticipant, err := srv.participantService.IsConferenceParticipant(int(req.UserId), int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	isOrganizer, err := srv.organizerService.IsConferenceOrganizer(int(req.UserId), int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.IsUserAssociatedWithConferenceResponse{
		IsAssociated: isParticipant || isOrganizer,
	}, nil
}
func (srv *ConferenceServ) GetOrganizerMetrics(ctx context.Context, req *pb.GetOrganizerMetricsRequest) (*pb.GetOrganizerMetricsResponse, error) {
	m, err := srv.organizerService.GetOrganizerMetrics(int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetOrganizerMetricsResponse{
		RunningConferences:        int32(m.RunningConferences),
		ParticipantsCount:         int32(m.ParticipantsCount),
		AverageParticipantsCount:  float32(m.AverageParticipantsCount),
		TotalOrganizedConferences: int32(m.TotalOrganizedConferences),
	}, nil
}
func (srv *ConferenceServ) GetParticipantsTrend(ctx context.Context, req *pb.ParticipantsTrendRequest) (*pb.ParticipantsTrendResponse, error) {
	trend, err := srv.organizerService.GetParticipantsJoiningTrend(int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	trendData := make([]*pb.ParticipantsTrendEntry, len(trend.Entries))
	for i, data := range trend.Entries {
		trendData[i] = &pb.ParticipantsTrendEntry{
			Date:            timestamppb.New(data.Date),
			NewParticipants: int32(data.NewParticipants),
		}
	}

	return &pb.ParticipantsTrendResponse{
		Trend: trendData,
	}, nil
}
func (srv *ConferenceServ) GetTickets(ctx context.Context, req *pb.TicketsRequest) (*pb.TicketsResponse, error) {
	t, m, err := srv.participantService.GetParticipantsTickets(int(req.UserId), commonService.Page{
		PageNumber: int(req.Page.Number),
		PageSize:   int(req.Page.Size),
	})
	log.Println(t)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tickets := make([]*pb.Ticket, len(t))
	log.Println(tickets)
	for i, ticket := range t {
		tickets[i] = &pb.Ticket{
			TicketId:     ticket.Id,
			ConferenceId: int32(ticket.ConferenceId),
		}
	}
	return &pb.TicketsResponse{
		TicketsPage: &pb.TicketsPage{
			Data: tickets,
			Meta: &pb.TicketsPageMeta{
				Page: &commonPb.PageInfo{
					Number:     int32(m.PageNumber),
					Size:       int32(m.PageSize),
					TotalItems: int32(m.TotalItems),
					TotalPages: int32(m.TotalPages),
				},
			},
		},
	}, nil
}

func (srv *ConferenceServ) GetConferenceOrganizer(ctx context.Context, req *pb.GetConferenceOrganizerRequest) (*pb.GetConferenceOrganizerResponse, error) {
	o, err := srv.organizerService.GetConferenceOrganizerId(int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetConferenceOrganizerResponse{
		OrganizerId: int32(o),
	}, nil
}
func (srv *ConferenceServ) ValidateTicket(ctx context.Context, req *pb.ValidateTicketRequest) (*pb.ValidateTicketResponse, error) {
	isValid, err := srv.participantService.IsTicketValid(req.TicketId, int(req.ConferenceId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ValidateTicketResponse{
		IsValid: isValid,
	}, nil
}
