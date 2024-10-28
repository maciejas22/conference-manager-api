package repositories

type Section struct {
	Id               int     `json:"id" db:"id"`
	TermsOfServiceId int     `json:"terms_of_service_id" db:"terms_of_service_id"`
	Title            string  `json:"title" db:"title"`
	Content          *string `json:"content" db:"content"`
	CreatedAt        string  `json:"created_at" db:"created_at"`
	JoinedAt         string  `json:"joined_at" db:"joined_at"`
	Subsections      []*Subsection
}

func (s *Section) TableName() string {
	return "public.sections"
}
