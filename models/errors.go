package models

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized user")
	ErrInvalidUser     = errors.New("invalid user")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidUserType = errors.New("invalid user type")
	ErrWrongPassword   = errors.New("wrong password")
	ErrUnexpected      = errors.New("unexpected error occured")
	ErrUserNotAdmin    = errors.New("only for admin user")
	ErrBookNotExists   = errors.New("book not exists")
)
