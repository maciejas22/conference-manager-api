package repository

import "time"

type Subsection struct {
	Id        int       `json:"id" db:"id"`
	SectionId int       `json:"section_id" db:"section_id"`
	Title     string    `json:"title" db:"title"`
	Content   *string   `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
