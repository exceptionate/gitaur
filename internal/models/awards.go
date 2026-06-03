package models

type Award struct {
	ID          int
	Title       string
	Issuer      string
	Type        string
	Tags        []string
	Description string
	Date        string
}

type Awards []Award
