package auth

import (
	"github.com/OhMinsSup/lafu-server/middlewares"
	"github.com/labstack/echo/v4"
)

// ApplyRoutes 라우터
func ApplyRoutes(e *echo.Group) {
	auth := e.Group("/auth")

	auth.POST("/register/local", localRegister)
	auth.POST("/login/local", localLogin)
	auth.POST("/logout", logout)

	auth.GET("/info", authInfo, middlewares.Authorized)
}
