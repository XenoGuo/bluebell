package mysql

import "errors"

var (
	ErrNotFound  = errors.New("dao: not found")
	ErrDuplicate = errors.New("dao: duplicate")
)
