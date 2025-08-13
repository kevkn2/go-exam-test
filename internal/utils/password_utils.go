package utils

import "golang.org/x/crypto/bcrypt"

type ComparePasswords struct {
	Password       string
	HashedPassword string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePasswordFail(cp ComparePasswords) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(cp.HashedPassword),
		[]byte(cp.Password),
	) != nil
}
