package util

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	key []byte = []byte("qiangjiyi@163.com")
)

// GenerateToken create new json web token
func GenerateToken() string {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + 1000),
		Issuer:    "qiangjiyi",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return ss
}

// VerifyToken check token if valid
func VerifyToken(token string) bool {
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parse with claims failed.", err)
		return false
	}
	return true
}
