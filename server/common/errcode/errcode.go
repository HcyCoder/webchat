package errcode

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrTokenExpired      = errors.New("token expired")
	ErrGroupNotFound     = errors.New("group not found")
	ErrNotGroupMember    = errors.New("not a group member")
	ErrMessageNotFound   = errors.New("message not found")
)
