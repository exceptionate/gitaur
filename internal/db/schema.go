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
	tags TEXT,
	desc TEXT,
	text TEXT,
	startDate TEXT,
	endDate TEXT
)
`

const AwardsSchema = `
CREATE TABLE IF NOT EXISTS awards (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	issuer TEXT NOT NULL,
	type TEXT NOT NULL,
	tags TEXT,
	description TEXT,
	date TEXT
)
`

const CertificationsSchema = `
CREATE TABLE IF NOT EXISTS certifications (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	platform TEXT NOT NULL,
	description TEXT,
	url TEXT,
	date TEXT,
	tags TEXT
)
`
