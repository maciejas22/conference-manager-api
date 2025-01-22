package repository

type TableName string

const (
	News           TableName = "public.news"
	Section        TableName = "public.sections"
	Subsection     TableName = "public.subsections"
	TermsOfService TableName = "public.terms_of_service"
)

func GetTableName(tableName TableName) string {
	return string(tableName)
}
