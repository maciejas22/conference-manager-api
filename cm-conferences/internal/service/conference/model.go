package service

import "time"

type Conference struct {
	Id                   int
	Title                string
	StartDate            time.Time
	EndDate              time.Time
	Location             string
	Website              *string
	Acronym              *string
	AdditionalInfo       *string
	ParticipantsLimit    *int
	RegistrationDeadline *time.Time
	TicketPrice          int
}

type ConferencesFilters struct {
	Title          *string
	AssociatedOnly *bool
	RunningOnly    *bool
}

type ConferencesMetrics struct {
	RunningConferences        int
	StartingInLessThan24Hours int
	TotalConducted            int
	ParticipantsToday         int
}

type AgendaItem struct {
	Id           int
	ConferenceId int
	StartTime    time.Time
	EndTime      time.Time
	Event        string
	Speaker      string
}

type UploadFileInput struct {
	FileName   string
	Base64File string
}

type CreateConferenceInput struct {
	Title                string
	StartDate            time.Time
	EndDate              time.Time
	Location             string
	Website              *string
	Acronym              *string
	AdditionalInfo       *string
	ParticipantsLimit    *int
	RegistrationDeadline *time.Time
	Agenda               []*AgendaItem
	Files                []*UploadFileInput
	TicketPrice          int
}

type DeleteFileInput struct {
	Key string
}

type ModifyConferenceFilesInput struct {
	Uploads []*UploadFileInput
	Deletes []*DeleteFileInput
}

type ModifyConferenceInput struct {
	ID                   int
	Title                *string
	StartDate            *time.Time
	EndDate              *time.Time
	Location             *string
	Website              *string
	Acronym              *string
	AdditionalInfo       *string
	ParticipantsLimit    *int
	RegistrationDeadline *time.Time
	Agenda               []*AgendaItem
	Files                *ModifyConferenceFilesInput
	TicketPrice          *int
}
