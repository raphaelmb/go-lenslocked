package models

import "database/sql"

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When look up a session this will be left empty.
	// Only store the hash of a session token in the database and cannot reverse it into raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}
