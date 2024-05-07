package logic

import (
	"backend/models"
	"backend/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取authorization header
		tokenString := c.GetHeader("Authorization")

		//验证token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort()
			return
		}

		//验证通过后获取claims中的Id
		userId := claims.UserId
		var user models.User
		mysql.DB.First(&user, userId) //默认用户存在
		c.Set("user", user)
		c.Next()
	}
}
