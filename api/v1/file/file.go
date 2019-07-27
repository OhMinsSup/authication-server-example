package file

import (
"github.com/labstack/echo/v4"
)

// ApplyRoutes 라우터
func ApplyRoutes(e *echo.Group) {
	file := e.Group("/file")

	file.POST("/create-url", createUrl)
	file.POST("/resize", resize)
}
