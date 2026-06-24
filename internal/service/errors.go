package service

import "errors"

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrMarshalJSON       = errors.New("invalid JSON in config")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrAlreadyVoted      = errors.New("user already voted in poll")
	ErrNotUserPoll       = errors.New("invalid user poll")
	ErrAuthOnly          = errors.New("this poll for only auth users")
	ErrInvalidUsername   = errors.New("username can only contain latin letters, numbers, and underscores")
	ErrMustContainLetter = errors.New("username must contain at least one latin letter")
	ErrDoubleUnderscores = errors.New("username cannot contain consecutive underscores")
	ErrUsernameTooBig    = errors.New("username exceeds maximum allowed length")
	ErrUsernameTooSmall  = errors.New("username is shorter than minimum allowed length")
	ErrPasswordTooBig    = errors.New("password exceeds maximum allowed length")
	ErrPasswordTooSmall  = errors.New("password is shorter than minimum allowed length")
	ErrAnswersNotCompare = errors.New("invalid vote request, answers length not compare with poll")
	ErrQuestionType      = errors.New("invalid vote request, answers not compare with question type")
	ErrQuestionNotFound  = errors.New("invalid vote request, question for answer not found")
	ErrOptionNotFound    = errors.New("invalid vote request, option for answer not found")
	ErrEmptyAnswer       = errors.New("invalid vote request, empty text type answer")
	ErrDuplicateOptions  = errors.New("invalid vote request, found duplicate of answer options")
)
