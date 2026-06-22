package service

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrMarshalJSON     = errors.New("invalid JSON in config")
	ErrInvalidToken    = errors.New("invalid or expired token")
	ErrAlreadyVoted    = errors.New("user already voted in poll")
	ErrNotUserPoll     = errors.New("invalid user poll")
	ErrAuthOnly        = errors.New("this poll for only auth users")
)
