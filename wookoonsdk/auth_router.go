package wookoonsdk

func (s *Server) useAuthRoutes() {
	sv := s.newAuthService()
	hd := s.newAuthHandler(sv)

	// API 路由组
	api := s.gin.Group("/api/cas")
	{
		// 公开接口
		api.GET("/login", hd.GetLoginURL)
		api.GET("/login-auto", hd.RedirectLogin)
		api.POST("/callback", hd.HandleCallback)

		// 需要认证的接口
		auth := api.Group("")
		auth.Use(s.AuthRequired())
		{
			auth.GET("/me", hd.GetMe)
			api.GET("/me-auto", hd.RedirectProfile)
			auth.POST("/logout", hd.Logout)
		}
	}
}
