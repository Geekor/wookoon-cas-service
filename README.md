# wookoon-cas-service (Go 版本 SDK)

## 代码架构
```
wookoonsdk/
├── server.go           # 服务入口
├── config.go           # 配置加载
├── cors.go             # CORS 中间件
├── casdoor_client.go   # Casdoor 客户端封装
├── auth_handler.go     # HTTP 处理器
├── auth_model_user.go  # 数据模型（防腐层）
├── auth_middleware.go  # JWT 验证中间件
├── auth_service.go     # 业务逻辑层
├── auth_router.go      # 路由注册
└── jwtutil.go          # JWT 签发/验证
```

## SDK 使用示例

```go
package main

import (
	"log"

	"github.com/geekor/wookoon-cas-service/wookoonsdk"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	s := wookoonsdk.NewServer(r)

	// 追加自定义接口：健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	if err := s.Run(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
```

运行

```bash
# 配置环境
cp .env.example .env

# 启动 DEV
go run main.go
```

## 前端接口说明

- `/api/cas/login` 登录(返回登录页 URL，需要前端解析并跳转)
- `/api/cas/login-auto` 登录(自动重定向到登录页面)
- `/api/cas/callback` 前端获取到登录 code + state 后，调用本接口换取 token
- `/api/cas/me` 获取当前用户信息
- `/api/cas/logout` 登出