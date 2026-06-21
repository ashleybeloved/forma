package service

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrMarshalJSON     = errors.New("invalid JSON in config")
	ErrInvalidToken    = errors.New("invalid or expired token")
)
