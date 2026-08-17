package wookoonsdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler 创建认证处理器
func (s *Server) newAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GetLoginURL 获取登录 URL
// GET /api/cas/login?redirect_uri=xxx
func (h *AuthHandler) GetLoginURL(c *gin.Context) {
	redirectURI := c.Query("redirect_uri")

	resp, err := h.authService.GetLoginURL(redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RedirectLogin 直接重定向到登录页（可选）
// GET /api/cas/login-auto?redirect_uri=xxx
func (h *AuthHandler) RedirectLogin(c *gin.Context) {
	redirectURI := c.Query("redirect_uri")

	resp, err := h.authService.GetLoginURL(redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, resp.LoginURL)
}

// HandleCallback 处理登录回调
// POST /api/cas/callback
func (h *AuthHandler) HandleCallback(c *gin.Context) {
	var req CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	resp, err := h.authService.HandleCallback(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "登录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RedirectProfile 直接重定向到个人资料页
// GET /api/cas/profile-auto
func (h *AuthHandler) RedirectProfile(c *gin.Context) {
	resp, err := h.authService.GetMyProfile()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, resp.ProfileURL)
}

// GetMe 获取当前用户信息
// GET /api/cas/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	// 从中间件注入的 Context 中获取用户信息
	user := &StandardUser{
		ID:          c.GetString("userId"),
		Username:    c.GetString("username"),
		DisplayName: c.GetString("displayName"),
		Email:       c.GetString("email"),
		Roles:       c.GetStringSlice("roles"),
	}

	c.JSON(http.StatusOK, user)
}

// Logout 登出
// POST /api/cas/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// 由于 JWT 是无状态的，登出主要由前端清除 Token 实现
	// 如果需要服务端登出，可以维护一个 Token 黑名单（Redis）

	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}
