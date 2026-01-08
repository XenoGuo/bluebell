package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int // 建议 12~14；先用 12 比较通用
}

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = 12
	}
	return &BcryptHasher{Cost: cost}
}

// Hash 生成可存库的 hash（包含 salt），直接存这个字符串即可
func (h *BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.Cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(b), nil
}

// Verify 校验密码（返回 nil 表示通过）
func (h *BcryptHasher) Verify(hash, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err == nil {
		return nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidCredentials
	}
	return fmt.Errorf("bcrypt verify: %w", err)
}

var ErrInvalidCredentials = fmt.Errorf("invalid credentials")
