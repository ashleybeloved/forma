package repository

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("user with this username already exists")
	ErrUserNotFound          = errors.New("user not found")
)

var (
	ErrPollNotFound  = errors.New("poll not found")
	ErrPollsNotFound = errors.New("polls not found")
)
