package models

type User struct {
	UserID   int64  `db:"user_id"`
	Gender   int32  `db:"gender"`
	Username string `db:"username"`
	Password string `db:"password"`
}
