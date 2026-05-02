package service

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")

	ErrCommentNotFound = errors.New("comment not found")
)
