package models

type Experience struct {
	ID          int
	Title       string
	Company     string
	StartDate   string
	EndDate     string
	Description string
	Tags        []string
}

type Experiences []Experience
