package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"shop_api/user_web/global"
	"shop_api/user_web/models"
	"time"
)

// JWTAuth 是一个JWT身份验证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取x-token信息
		token := c.Request.Header.Get("x-token")
		zap.S().Debug("token-x的值为:" + token)

		if token == "" {
			// 如果令牌为空，则返回未登录的错误响应
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// 解析令牌
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				// 如果令牌过期，则返回授权已过期的错误响应
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}

			// 令牌无效或解析错误，则返回未登录的错误响应
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}

		// 将claims设置到gin的上下文中，以便后续处理程序可以使用
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

// JWT 是JWT结构体
type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

// NewJWT 创建一个JWT实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.Key), // 把全局密钥作为SigningKey
	}
}

// CreateToken 创建一个JWT令牌
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	// 使用SigningKey对claims进行签名，生成令牌字符串
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析JWT令牌
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 解析令牌字符串，并使用SigningKey验证令牌的合法性
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		// 处理解析过程中可能出现的错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// 令牌已过期
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	// 验证令牌是否有效以及是否属于自定义声明类型，如果有效则返回claims
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// RefreshToken 刷新JWT令牌
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 设置jwt.TimeFunc以便在解析过期令牌时时间一直为0
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}

	// 如果令牌有效且属于自定义声明类型，则刷新令牌的过期时间，并生成新的令牌字符串
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		// 将jwt.TimeFunc设置回默认值，以便生成具有正确过期时间的新令牌
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}

	return "", TokenInvalid
}
