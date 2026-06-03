package db

const UserSchema = `
CREATE TABLE IF NOT EXISTS user (
	firstName TEXT NOT NULL,
	lastName TEXT NOT NULL,
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	avatar TEXT,
	linkedin TEXT,
	portfolio TEXT,
	leetcode TEXT,
	phoneNumber TEXT
)
`

const ProjectsSchema = `
CREATE TABLE IF NOT EXISTS projects (
	name TEXT NOT NULL,
	repo TEXT NOT NULL,
	owner TEXT NOT NULL,
	url TEXT NOT NULL,
	tech TEXT,
	desc TEXT,
	text TEXT,
	startDate TEXT,
	endDate TEXT
)
`
