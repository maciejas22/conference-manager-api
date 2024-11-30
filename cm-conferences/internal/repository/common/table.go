package common

type TableName string

const (
	AgendaTable                TableName = "public.agenda"
	ConferenceTable            TableName = "public.conferences"
	ConferenceOrganizerTable   TableName = "public.conference_organizers"
	ConferenceParticipantTable TableName = "public.conference_participants"
	UserTable                  TableName = "public.users"
)

func GetTableName(tableName TableName) string {
	return string(tableName)
}
