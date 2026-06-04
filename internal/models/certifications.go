package models

type Certification struct {
	ID          int
	Title       string
	Platform    string
	Description string
	Url         string
	Date        string
	Tags        []string
}

type Certifications []Certification
