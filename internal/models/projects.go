package models

type Project struct {
	Name      string
	Repo      string
	Owner     string
	Url       string
	Tech      []string
	Desc      string
	Text      string
	StartDate string
	EndDate   string
}

type Projects []Project
