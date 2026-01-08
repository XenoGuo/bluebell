package mysql

import (
	"bluebell/models"
	"bluebell/pkg/auth"
	"database/sql"
	"errors"
	"fmt"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	sqlStr := `insert into user (user_id,username,password) values (?,?,?)`
	hasher := auth.NewBcryptHasher(12)
	hash, err := hasher.Hash(user.Password)
	if err != nil {
		return err
	}
	_, err = db.Exec(sqlStr, user.UserID, user.Username, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDuplicate
		}
		return err
	}
	return err
}

func CheckUserExist(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("CheckUserExist err: %w", err)
	}
	return count > 0, nil
}

// 可复用模板
func SelectUser(username string) (user *models.User, err error) {
	u := new(models.User)
	sqlStr := `select user_id,username,password from user where username=?`
	if err := db.Get(u, sqlStr, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("SelectUser err: %w", err)
	}
	return u, nil
}

func GetUserById(id int64) (*models.User, error) {
	u := new(models.User)
	sqlStr := `select user_id,username,gender from user where user_id=?`
	err := db.Get(u, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("GetUserById err: %w", err)
	}
	return u, nil
}
