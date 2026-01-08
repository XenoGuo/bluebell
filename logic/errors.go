package logic

import "errors"

var (
	ErrCommunityNotFound = errors.New("社区不存在")
	ErrPostNotFound      = errors.New("帖子不存在")
	ErrInvalidPassword   = errors.New("用户名或密码错误")
	ErrUserNotExists     = errors.New("用户不存在")
	ErrUserExists        = errors.New("用户已存在")
)
