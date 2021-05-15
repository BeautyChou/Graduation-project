package Middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const tokenExpireDuration = time.Hour * 3

var mySecret = []byte("这是我的密钥")

type Myclaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Cors 解决跨域问题的中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")
		if len(origin) == 0 {
			origin = c.Request.Header.Get("Origin")
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Token, user-agent")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, DELETE, PUT")

		if c.Copy().Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			return
		} else {
			c.Next()
		}
	}
}

//Json-web-token
//	GenToken   创建令牌
//	ParseToken 解析令牌
//  JWTAuthMiddleWare 解决认证请求的中间件，外部依赖于GenToken产生的令牌，内部依赖于ParseToken提供的解析令牌的功能
func GenToken(username string) (string, error) {
	claims := Myclaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
			Issuer:    "ZhouChaoYiProject",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*Myclaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Myclaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Myclaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuthMiddleWare() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Token")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"msg": "请求头中Token为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "8a54sh") {
			c.JSON(http.StatusOK, gin.H{
				"msg": "Token格式错误",
			})
			c.Abort()
			return
		}
		_, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "Token无效",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
