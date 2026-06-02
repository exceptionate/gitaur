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
	githubToken TEXT,
	phoneNumber TEXT
)
`
