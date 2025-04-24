package logic

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("xxxxxxx")

// type Claims struct {
// 	UserId uint
// 	jwt.StandardClaims
// }

//	func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
//		claims := &Claims{}
//		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
//			return jwtKey, nil
//		})
//		return token, claims, err
//	}
type Claims struct {
	LoginTime time.Time
	UserId    uint
	Identity  string
	PU_uid    uint
	jwt.StandardClaims
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

func GenToken(id uint, identity string, pu_uid uint) (string, error) {
	claims := Claims{
		LoginTime: time.Now(),
		UserId:    id,
		Identity:  identity,
		PU_uid:    pu_uid,
	}
	// claims.ExpiresAt = time.Now().Add(3 * time.Hour).Unix()
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtKey)
}

// func GetToken(userid uint) (string, error) {
// 	claims := Claims{
// 		UserId: userid,
// 	}
// 	//尝试获取一个token
// 	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return tokenClaims.SignedString(jwtKey)
// }
