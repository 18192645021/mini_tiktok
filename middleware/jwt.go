package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mock_douyin_project/dao"
	"net/http"
	"time"
)

var jwtKey = []byte("acking-you.xyz")

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// GenerateToken 登录成功后，传入userId生成token
func GenerateToken(user dao.UserLogin) (string, error) {
	// 设置过期时间为2小时
	expirationTime := time.Now().Add(24 * 7 * time.Hour)
	claims := &Claims{
		UserId: user.UserInfoId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "happy_boy_111",
			Subject:   "L_Q__",
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析jwt
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		// token不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, dao.CommonResponse{StatusCode: 401, StatusMsg: "用户未登录"})
			c.Abort()
			return
		}
		// 解析token
		token, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, dao.CommonResponse{StatusCode: 403, StatusMsg: err.Error()})
			c.Abort()
			return
		}
		//// token超时检验
		//if time.Now().Unix() > token.ExpiresAt {
		//	c.JSON(http.StatusOK, dao.CommonResponse{
		//		StatusCode: 404, StatusMsg: "token过期",
		//	})
		//}
		// 将当前请求的userId信息保存到请求的上下文c上
		c.Set("user_id", token.UserId)
		c.Next() // 后续的处理函数可以用过c.Get("user_id")来获取当前请求的用户信息
	}
}
