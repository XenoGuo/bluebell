package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/auth"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"errors"
)

// 存放业务逻辑的代码
func SignUp(p *models.ParamSignUp) (err error) {
	// 构造一个User实例
	u := &models.User{
		Username: p.Username,
		Password: p.Password,
		UserID:   snowflake.GenID(),
	}
	// 入库
	if err := mysql.InsertUser(u); err != nil {
		if errors.Is(err, mysql.ErrDuplicate) {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func Login(p *models.ParamLogin) (token string, err error) {
	// 查用户
	user, err := mysql.SelectUser(p.Username)
	if err != nil {
		if errors.Is(err, mysql.ErrNotFound) {
			return "", ErrUserNotExists
		}
		return "", err
	}
	// 参数校验
	hasher := auth.NewBcryptHasher(12)

	if err = hasher.Verify(user.Password, p.Password); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return "", ErrInvalidPassword
		}
		return "", err
	}
	// 生成JWT
	return jwt.GenToken(user.UserID, user.Username)
}
