package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JWTToken jwtToken

type jwtToken struct{}

// CustomClaims 定义token 携带的信息
type CustomClaims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	*jwt.StandardClaims
}

// SECRET 加解密因子
const (
	UserSecret          = "devops"       //普通用户加密因子
	UserExpireDuration  = time.Hour * 2  //普通用户过期时间2h
	AdminSecret         = "admin_devops" //管理员用户加密因子
	AdminExpireDuration = time.Hour * 2  //管理员用户过期时间2H
)

// ParseToken 解析token
func (*jwtToken) ParseToken(tokenString string, secret string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		Logger.Error("Token parse failed" + err.Error())
		//处理 token 的各种报错
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid")
			}
		}
		//处理 token 的各种报错
		//if ve, ok := err.(*jwt.ValidationError); ok {
		//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
		//				return nil, errors.New("TokenMalformed")
		//			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
		//				return nil, errors.New("TokenExpired")
		//			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
		//				return nil, errors.New("TokenNotValidYet")
		//			} else {
		//				return nil, errors.New("TokenInvalid")
		//			}
		//		}
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("解析Token失败")
}

// CreateJwtToken 后端产生token 用于验证身份
func CreateJwtToken(id int, username string, expireDuration time.Duration, secret string) (string, error) {
	// 自定义claims	jwt第一段加密
	claims := &CustomClaims{
		id,
		username,
		&jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,              // token信息生效时间
			ExpiresAt: time.Now().Add(expireDuration).Unix(), // 过期时间
			Issuer:    "kube_ops",                            // 发布者
		},
	}
	// 对自定义claims加密,jwt.SigningMethodHS256是加密算法得到第二部分
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 给这个token盐加密 第三部分,得到一个完整的三段的加密,真正意义上的加密
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
