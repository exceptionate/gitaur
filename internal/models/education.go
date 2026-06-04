package models

type Education struct {
	ID          int
	School      string
	Degree      string
	Field       string
	StartDate   string
	EndDate     string
	Grade       string
	Description string
}

type EducationList []Education
