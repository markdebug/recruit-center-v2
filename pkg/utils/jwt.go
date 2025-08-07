package utils

import (
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("your-secret-key") // 建议从配置文件读取

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
