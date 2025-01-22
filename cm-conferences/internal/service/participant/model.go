package service

type Ticket struct {
	Id           string
	ConferenceId int
}

type ConferenceParticipant struct {
	UserId       int
	ConferenceId int
	Ticket       Ticket
}

type ConferenceParticipantsCount struct {
	ConferenceId      int
	ParticipantsCount int
}
