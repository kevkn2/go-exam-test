package utils

import (
	"exam-test/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTUtils interface {
	GenerateToken(id uint) (string, error)
	ExtractToken(ctx *gin.Context) string
	TokenValid(ctx *gin.Context) (uint, error)
}

type jwtUtils struct {
	SecretKey     string
	ExpireMinutes time.Duration
}

// ExtractToken implements JWTUtils.
func (j *jwtUtils) ExtractToken(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer" {
		return token[7:]
	}
	return ""
}

// GenerateToken implements JWTUtils.
func (j *jwtUtils) GenerateToken(id uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(j.ExpireMinutes).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// TokenValid implements JWTUtils.
func (j *jwtUtils) TokenValid(ctx *gin.Context) (uint, error) {
	tokenStr := j.ExtractToken(ctx)

	if tokenStr == "" {
		return 0, fmt.Errorf("missing token")
	}

	token, err := jwt.Parse(
		tokenStr,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token")
	}

	userId := claims["user_id"].(float64)
	return uint(userId), nil
}

func NewJWTUtils(envConfig config.EnvConfig) JWTUtils {
	return &jwtUtils{
		SecretKey:     envConfig.SECRETKEY,
		ExpireMinutes: time.Duration(envConfig.ACCESS_TOKEN_EXPIRE_MINUTES) * time.Minute,
	}
}
