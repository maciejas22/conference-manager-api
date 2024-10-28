package repositories

type Subsection struct {
	Id        int     `json:"id" db:"id"`
	SectionId int     `json:"section_id" db:"section_id"`
	Title     string  `json:"title" db:"title"`
	Content   *string `json:"content" db:"content"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	JoinedAt  string  `json:"joined_at" db:"joined_at"`
}

func (s *Subsection) TableName() string {
	return "public.subsections"
}
