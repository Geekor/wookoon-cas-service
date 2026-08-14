package wookoonsdk

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*AppConfig
	*AuthConfig
	*JwtConfig
	*CorsConfig

	gin *gin.Engine
}

func NewServer(r *gin.Engine) *Server {
	s := &Server{
		gin: r,
	}

	// 加载配置
	s.loadConfigs()

	gin.SetMode(s.AppConfig.ServerMode)

	// 全局中间件
	s.useMiddlewareCORS()

	// 注册路由
	s.useAuthRoutes()

	return s
}

// 启动服务
func (s *Server) Run() error {
	addr := ":" + s.AppConfig.ServerPort
	log.Printf("🚀 Wookoon Auth Service 启动于 %s", addr)
	log.Printf("📍 Casdoor Endpoint: %s", s.AuthConfig.Endpoint)
	log.Printf("🔐 JWT Issuer: %s", s.JwtConfig.JWTIssuer)

	return s.gin.Run()
}
