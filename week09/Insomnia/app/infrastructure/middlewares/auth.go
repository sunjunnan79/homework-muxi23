package middlewares

import (
	"Insomnia/app/api/response"
	"Insomnia/app/infrastructure/jwt"
	"github.com/gin-gonic/gin"
)

// UseJwt 返回一个 gin.HandlerFunc，用于验证 JWT
func UseJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := jwt.Parse(token)
		if err != nil {
			response.FailMsg(c, "您还未登录!")
			c.Abort()
			return
		}
		c.Set("Uuid", claims["Uuid"])
		c.Next()
	}
}
