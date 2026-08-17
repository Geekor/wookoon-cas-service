package wookoonsdk

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired JWT 验证中间件
func (s *Server) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证信息",
			})
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "认证格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := JwtParseToken(tokenString, s.JwtConfig.JWTSecret)
		if err != nil {
			status := http.StatusUnauthorized
			message := "认证失败"
			if err == ErrTokenExpired {
				message = "Token 已过期"
			}

			c.JSON(status, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}

		// 将用户信息注入 Context，供后续 Handler 使用
		c.Set("tokenString", tokenString)
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("displayName", claims.DisplayName)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}
