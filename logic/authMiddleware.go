package logic

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// fmt.Println("触发中间件")
		//获取authorization header
		tokenString := c.GetHeader("Authorization")

		//验证token格式
		// if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		// 	c.Abort()
		// 	return
		// }
		// tokenString = tokenString[7:]

		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token无效"})
			c.Abort()
			return
		}

		//验证通过后获取claims中的Id
		userId := claims.UserId
		identity := claims.Identity
		pu_uid := claims.PU_uid

		c.Set("userId", userId)
		c.Set("identity", identity)
		c.Set("pu_uid", pu_uid)
		c.Next()
	}
}
