package auth

import (
"github.com/labstack/echo/v4"
)

// ApplyRoutes 라우터
func ApplyRoutes(e *echo.Group) {
	auth := e.Group("/auth")

	auth.POST("/register/local", localRegister)
	auth.POST("/login/local", localLogin)
	auth.POST("/logout", logout)
	auth.GET("/test", test)
}
