package models

import (
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims 自定义的JWT声明结构体
type CustomClaims struct {
	ID                 uint   // 用户ID
	Name               string // 用户姓名
	AuthorityId        uint   // 用户权限ID
	jwt.StandardClaims        // JWT标准声明
}
