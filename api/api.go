package api

import (
	v1 "github.com/OhMinsSup/lafu-server/api/v1"
	_ "github.com/OhMinsSup/lafu-server/docs"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

// ApplyRoutes 라우터
func ApplyRoutes(e *echo.Echo) {
	url := echoSwagger.URL("http://localhost:7000/api/swagger/doc.json")
	api := e.Group("/api")

	api.GET("/swagger/*", echoSwagger.EchoWrapHandler(url))
	v1.ApplyRoutes(api)
}
