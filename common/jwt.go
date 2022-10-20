// @Title  jwt
// @Description  该文件用于提供加密为token和解析token的函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package common

import (
	"MGA_OJ/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

// Claims			定义了user的token
type Claims struct {
	UserId             uint // 用户的id
	jwt.StandardClaims      // usertoken
}

// @title    ReleaseToken
// @description   生成用户的token
// @auth      MGAronya（张健）       2022-9-16 12:07
// @param     user model.User       接收一个用户
// @return    string, error         返回该用户的token，或者返回error
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "MGAronya",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	// TODO err不为空，则返回错误
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// @title    ParseToken
// @description   对token进行解析
// @auth      MGAronya（张健）       2022-9-16 12:07
// @param    tokenString string       接收一个token
// @return   *jwt.Token, *Claims, error         返回token中包含的信息
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	// TODO 解析token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
