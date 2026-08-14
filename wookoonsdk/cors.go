package wookoonsdk

import (
	"time"

	"github.com/gin-contrib/cors"
)

// CORS 跨域中间件
func (s *Server) useMiddlewareCORS() *Server {
	f := cors.New(cors.Config{
		AllowOrigins:     s.CorsConfig.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	s.gin.Use(f)

	return s
}
