package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/maciejas22/conference-manager-api/cm-gateway/config"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/graph"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/graph/model"
	middlewares "github.com/maciejas22/conference-manager-api/cm-gateway/internal/middleware"
	"github.com/maciejas22/conference-manager-api/cm-gateway/internal/util"
	authPb "github.com/maciejas22/conference-manager-api/cm-proto/auth"
	"github.com/maciejas22/conference-manager-api/cm-proto/common"
	pb "github.com/maciejas22/conference-manager-api/cm-proto/conference"
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Agenda is the resolver for the agenda field.
func (r *conferenceResolver) Agenda(ctx context.Context, obj *model.Conference) ([]*model.AgendaItem, error) {
	a, err := r.conferenceServiceClient.GetConferenceAgenda(ctx, &pb.AgendaRequest{
		ConferenceId: obj.ID,
	})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	agendaItems := make([]*model.AgendaItem, len(a.Agenda))
	for i, agendaItem := range a.Agenda {
		agendaItems[i] = &model.AgendaItem{
			ID:        agendaItem.Id,
			Event:     agendaItem.Event,
			Speaker:   agendaItem.Speaker,
			StartTime: agendaItem.StartTime.AsTime(),
			EndTime:   agendaItem.EndTime.AsTime(),
		}
	}

	return agendaItems, nil
}

// Files is the resolver for the files field.
func (r *conferenceResolver) Files(ctx context.Context, obj *model.Conference) ([]*model.File, error) {
	files, err := r.s3Client.GetFilesFromFolder(ctx, config.AppConfig.Buckets.ConferenceFiles, strconv.Itoa(int(obj.ID)))
	if err != nil {
		return nil, err
	}
	var parsedFiles []*model.File
	for _, f := range files {
		parsedFiles = append(parsedFiles, &model.File{
			Key:  f.Key,
			URL:  f.URL,
			Size: int32(f.Size),
		})
	}

	return parsedFiles, nil
}

// Metrics is the resolver for the metrics field.
func (r *conferencesPageResolver) Metrics(ctx context.Context, obj *model.ConferencesPage) (*model.ConferencesMetrics, error) {
	m, err := r.conferenceServiceClient.GetConferencesMetrics(ctx, &pb.ConferencesMetricsRequest{})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	return &model.ConferencesMetrics{
		RunningConferences:        m.Metrics.RunningConferences,
		StartingInLessThan24Hours: m.Metrics.StartingInLessThan24Hours,
		TotalConducted:            m.Metrics.TotalConducted,
		ParticipantsToday:         m.Metrics.ParticipantsToday,
	}, nil
}

// CreateConference is the resolver for the createConference field.
func (r *mutationResolver) CreateConference(ctx context.Context, createConferenceInput model.CreateConferenceInput) (int32, error) {
	si := middlewares.GetSessionInfo(ctx)

	var website, acronym, additionalInfo string
	if createConferenceInput.Website != nil {
		website = *createConferenceInput.Website
	}
	if createConferenceInput.Acronym != nil {
		acronym = *createConferenceInput.Acronym
	}
	if createConferenceInput.AdditionalInfo != nil {
		additionalInfo = *createConferenceInput.AdditionalInfo
	}

	var participantsLimit, ticketPrice int32
	if createConferenceInput.ParticipantsLimit != nil {
		participantsLimit = *createConferenceInput.ParticipantsLimit
	}
	if createConferenceInput.TicketPrice != nil {
		ticketPrice = *createConferenceInput.TicketPrice
	}

	var registrationDeadline *timestamppb.Timestamp = nil
	if createConferenceInput.RegistrationDeadline != nil {
		registrationDeadline = timestamppb.New(*createConferenceInput.RegistrationDeadline)
	}

	agendaItems := make([]*pb.CreateAgendaItemInput, len(createConferenceInput.Agenda))
	for i, agendaItem := range createConferenceInput.Agenda {
		agendaItems[i] = &pb.CreateAgendaItemInput{
			Event:     agendaItem.CreateItem.Event,
			Speaker:   agendaItem.CreateItem.Speaker,
			StartTime: timestamppb.New(agendaItem.CreateItem.StartTime),
			EndTime:   timestamppb.New(agendaItem.CreateItem.EndTime),
		}
	}

	cId, err := r.conferenceServiceClient.CreateConference(ctx, &pb.CreateConferenceRequest{
		UserId: si.UserId,
		CreateConferenceInput: &pb.CreateConferenceInput{
			Title:                createConferenceInput.Title,
			StartDate:            timestamppb.New(createConferenceInput.StartDate),
			EndDate:              timestamppb.New(createConferenceInput.EndDate),
			Location:             createConferenceInput.Location,
			Website:              website,
			Acronym:              acronym,
			AdditionalInfo:       additionalInfo,
			ParticipantsLimit:    participantsLimit,
			RegistrationDeadline: registrationDeadline,
			TicketPrice:          ticketPrice,
			Agenda:               agendaItems,
		},
	})
	if err != nil {
		return 0, gqlerror.Errorf(err.Error())
	}

	for _, f := range createConferenceInput.Files {
		key := strconv.Itoa(int(cId.ConferenceId)) + "/" + f.UploadFile.FileName
		file, err := util.FileConvert(f.UploadFile.Base64Content)
		if err != nil {
			return 0, err
		}

		err = r.s3Client.UploadFile(ctx, config.AppConfig.Buckets.ConferenceFiles, key, file)
		if err != nil {
			return 0, err
		}
	}

	return cId.ConferenceId, nil
}

// ModifyConference is the resolver for the modifyConference field.
func (r *mutationResolver) ModifyConference(ctx context.Context, input model.ModifyConferenceInput) (int32, error) {
	var title, location, website, acronym, additionalInfo string
	if input.Title != nil {
		title = *input.Title
	}
	if input.Location != nil {
		location = *input.Location
	}
	if input.Website != nil {
		website = *input.Website
	}
	if input.Acronym != nil {
		acronym = *input.Acronym
	}
	if input.AdditionalInfo != nil {
		additionalInfo = *input.AdditionalInfo
	}

	var participantsLimit, ticketPrice int32
	if input.ParticipantsLimit != nil {
		participantsLimit = *input.ParticipantsLimit
	}
	if input.TicketPrice != nil {
		ticketPrice = *input.TicketPrice
	}

	var startDate, endDate *timestamppb.Timestamp
	var registrationDeadline *timestamppb.Timestamp = nil
	if input.StartDate != nil {
		startDate = timestamppb.New(*input.StartDate)
	}
	if input.EndDate != nil {
		endDate = timestamppb.New(*input.EndDate)
	}
	if input.RegistrationDeadline != nil {
		registrationDeadline = timestamppb.New(*input.RegistrationDeadline)
	}

	agendaItems := make([]*pb.ModifyAgendaItemInput, len(input.Agenda))
	for i, agendaItem := range input.Agenda {
		if agendaItem.CreateItem != nil {
			agendaItems[i] = &pb.ModifyAgendaItemInput{
				CreateItem: &pb.CreateAgendaItemInput{
					Event:     agendaItem.CreateItem.Event,
					Speaker:   agendaItem.CreateItem.Speaker,
					StartTime: timestamppb.New(agendaItem.CreateItem.StartTime),
					EndTime:   timestamppb.New(agendaItem.CreateItem.EndTime),
				},
			}
		} else {
			agendaItems[i] = &pb.ModifyAgendaItemInput{
				DeleteItem: *agendaItem.DeleteItem,
			}
		}
	}

	cId, err := r.conferenceServiceClient.ModifyConference(ctx, &pb.ModifyConferenceRequest{
		Input: &pb.ModifyConferenceInput{
			Id:                   input.ID,
			Title:                title,
			StartDate:            startDate,
			EndDate:              endDate,
			Location:             location,
			Website:              website,
			Acronym:              acronym,
			AdditionalInfo:       additionalInfo,
			ParticipantsLimit:    participantsLimit,
			RegistrationDeadline: registrationDeadline,
			TicketPrice:          ticketPrice,
			Agenda:               agendaItems,
		},
	})
	if err != nil {
		return 0, gqlerror.Errorf(err.Error())
	}

	for _, f := range input.Files {
		if f.DeleteFile != nil {
			err = r.s3Client.DeleteFile(ctx, config.AppConfig.Buckets.ConferenceFiles, f.DeleteFile.Key)
			if err != nil {
				return 0, err
			}
		}
		if f.UploadFile != nil {
			key := strconv.Itoa(int(cId.ConferenceId)) + "/" + f.UploadFile.FileName
			file, err := util.FileConvert(f.UploadFile.Base64Content)
			if err != nil {
				return 0, err
			}

			err = r.s3Client.UploadFile(ctx, config.AppConfig.Buckets.ConferenceFiles, key, file)
			if err != nil {
				return 0, err
			}
		}
	}

	return cId.ConferenceId, nil
}

// AddUserToConference is the resolver for the addUserToConference field.
func (r *mutationResolver) AddUserToConference(ctx context.Context, conferenceID int32) (*string, error) {
	si := middlewares.GetSessionInfo(ctx)
	c, err := r.conferenceServiceClient.GetConferences(ctx, &pb.ConferencesRequest{
		ConferenceIds: []int32{conferenceID},
	})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	conference := c.Conferences[0]
	if conference.TicketPrice == 0 {
		_, err := r.conferenceServiceClient.AddUserToConference(ctx, &pb.AddUserToConferenceRequest{
			ConferenceId: conferenceID,
			UserId:       si.UserId,
		})
		if err != nil {
			return nil, gqlerror.Errorf(err.Error())
		}
		return nil, nil
	}

	o, err := r.conferenceServiceClient.GetConferenceOrganizer(ctx, &pb.GetConferenceOrganizerRequest{
		ConferenceId: conferenceID,
	})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	organizerDetails, err := r.authServiceClient.UserProfileById(ctx, &authPb.UserProfileByIdRequest{
		UserId: o.OrganizerId,
	})

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(conference.TicketPrice)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		TransferData: &stripe.PaymentIntentTransferDataParams{
			Destination: stripe.String(organizerDetails.User.StripeAccountId),
		},
		Metadata: map[string]string{
			"conference_id": strconv.Itoa(int(conferenceID)),
			"user_id":       strconv.Itoa(int(si.UserId)),
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, errors.New("Error creating payment intent")
	}

	return stripe.String(pi.ClientSecret), nil
}

// RemoveUserFromConference is the resolver for the removeUserFromConference field.
func (r *mutationResolver) RemoveUserFromConference(ctx context.Context, conferenceID int32) (int32, error) {
	si := middlewares.GetSessionInfo(ctx)
	_, err := r.conferenceServiceClient.RemoveUserFromConference(ctx, &pb.RemoveUserFromConferenceRequest{
		ConferenceId: conferenceID,
		UserId:       si.UserId,
	})
	if err != nil {
		return 0, gqlerror.Errorf(err.Error())
	}

	return conferenceID, nil
}

func (r *mutationResolver) ValidateTicket(ctx context.Context, input model.ValidateTicketInput) (bool, error) {
	ticket, err := r.conferenceServiceClient.ValidateTicket(ctx, &pb.ValidateTicketRequest{
		ConferenceId: input.ConferenceID,
		TicketId:     input.TicketID,
	})
	if err != nil {
		return false, gqlerror.Errorf(err.Error())
	}

	return ticket.IsValid, nil
}

// Conference is the resolver for the conference field.
func (r *queryResolver) Conference(ctx context.Context, id int32) (*model.Conference, error) {
	c, err := r.conferenceServiceClient.GetConferences(ctx, &pb.ConferencesRequest{
		ConferenceIds: []int32{id},
	})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	var registrationDeadline *time.Time
	if c.Conferences[0].RegistrationDeadline != nil {
		date := c.Conferences[0].RegistrationDeadline.AsTime()
		registrationDeadline = &date
	}
	return &model.Conference{
		ID:                   c.Conferences[0].Id,
		Title:                c.Conferences[0].Title,
		StartDate:            c.Conferences[0].StartDate.AsTime(),
		EndDate:              c.Conferences[0].EndDate.AsTime(),
		Location:             c.Conferences[0].Location,
		Website:              &c.Conferences[0].Website,
		Acronym:              &c.Conferences[0].Acronym,
		AdditionalInfo:       &c.Conferences[0].AdditionalInfo,
		ParticipantsCount:    c.Conferences[0].ParticipantsCount,
		ParticipantsLimit:    &c.Conferences[0].ParticipantsLimit,
		RegistrationDeadline: registrationDeadline,
		EventsCount:          c.Conferences[0].EventsCount,
		TicketPrice:          &c.Conferences[0].TicketPrice,
	}, nil
}

// Conferences is the resolver for the conferences field.
func (r *queryResolver) Conferences(ctx context.Context, page *model.Page, sort *model.Sort, filters *model.ConferencesFilters) (*model.ConferencesPage, error) {
	si := middlewares.GetSessionInfo(ctx)
	var sortOrder common.Order
	if sort != nil && sort.Order == model.OrderAsc {
		sortOrder = common.Order_ASC
	} else {
		sortOrder = common.Order_DESC
	}
	cSort := &common.Sort{
		Column: "id",
		Order:  sortOrder,
	}

	cFilters := &pb.ConferencesFilters{}
	if filters != nil {
		if filters.AssociatedOnly != nil {
			cFilters.AssociatedOnly = *filters.AssociatedOnly
		}
		if filters.RunningOnly != nil {
			cFilters.RunningOnly = *filters.RunningOnly
		}
		if filters.Title != nil {
			cFilters.Title = *filters.Title
		}
	}

	cPage := &common.Page{
		Number: 1,
		Size:   10,
	}
	if page != nil {
		cPage.Number = page.Number
		cPage.Size = page.Size
	}

	cIds, err := r.conferenceServiceClient.GetConferencesPage(ctx, &pb.ConferencesPageRequest{
		UserId:  si.UserId,
		Page:    cPage,
		Sort:    cSort,
		Filters: cFilters,
	})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	c, err := r.conferenceServiceClient.GetConferences(ctx, &pb.ConferencesRequest{
		ConferenceIds: cIds.ConferencesPage.ConferenceIds,
	})
	if err != nil {
		return nil, gqlerror.Errorf(err.Error())
	}

	confMap := make(map[int32]*pb.Conference)
	for _, conference := range c.Conferences {
		confMap[conference.Id] = conference
	}

	orderedConferences := make([]*model.Conference, len(cIds.ConferencesPage.ConferenceIds))
	for i, conferenceId := range cIds.ConferencesPage.ConferenceIds {
		registrationDeadline := confMap[conferenceId].RegistrationDeadline.AsTime()
		orderedConferences[i] = &model.Conference{
			ID:                   confMap[conferenceId].Id,
			Title:                confMap[conferenceId].Title,
			StartDate:            confMap[conferenceId].StartDate.AsTime(),
			EndDate:              confMap[conferenceId].EndDate.AsTime(),
			Location:             confMap[conferenceId].Location,
			Website:              &confMap[conferenceId].Website,
			Acronym:              &confMap[conferenceId].Acronym,
			AdditionalInfo:       &confMap[conferenceId].AdditionalInfo,
			ParticipantsCount:    confMap[conferenceId].ParticipantsCount,
			ParticipantsLimit:    &confMap[conferenceId].ParticipantsLimit,
			RegistrationDeadline: &registrationDeadline,
			EventsCount:          confMap[conferenceId].EventsCount,
			TicketPrice:          &confMap[conferenceId].TicketPrice,
		}
	}

	return &model.ConferencesPage{
		Data: orderedConferences,
		Meta: &model.ConferenceMeta{
			Page: &model.PageInfo{
				Number:     cIds.ConferencesPage.Meta.Page.Number,
				Size:       cIds.ConferencesPage.Meta.Page.Size,
				TotalItems: cIds.ConferencesPage.Meta.Page.TotalItems,
				TotalPages: cIds.ConferencesPage.Meta.Page.TotalPages,
			},
		},
	}, nil
}

func (r *Resolver) Conference() graph.ConferenceResolver { return &conferenceResolver{r} }

func (r *Resolver) ConferencesPage() graph.ConferencesPageResolver {
	return &conferencesPageResolver{r}
}

func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type conferenceResolver struct{ *Resolver }
type conferencesPageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }