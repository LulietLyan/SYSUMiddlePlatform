package logic

import (
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

func GetToken(userid uint) (string, error) {
	claims := Claims{
		UserId: userid,
	}
	//尝试获取一个token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtKey)
}
